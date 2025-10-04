package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/lin-snow/ech0/internal/transaction"

	authModel "github.com/lin-snow/ech0/internal/model/auth"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/echo"
	commonRepository "github.com/lin-snow/ech0/internal/repository/common"
	repository "github.com/lin-snow/ech0/internal/repository/echo"
	commonService "github.com/lin-snow/ech0/internal/service/common"
	fediverseService "github.com/lin-snow/ech0/internal/service/fediverse"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

type EchoService struct {
	txManager        transaction.TransactionManager
	commonService    commonService.CommonServiceInterface
	echoRepository   repository.EchoRepositoryInterface
	commonRepository commonRepository.CommonRepositoryInterface
	fediverseService fediverseService.FediverseServiceInterface
}

func NewEchoService(
	tm transaction.TransactionManager,
	commonService commonService.CommonServiceInterface,
	echoRepository repository.EchoRepositoryInterface,
	commonRepository commonRepository.CommonRepositoryInterface,
	fediverseService fediverseService.FediverseServiceInterface,
) EchoServiceInterface {
	return &EchoService{
		txManager:        tm,
		commonService:    commonService,
		echoRepository:   echoRepository,
		commonRepository: commonRepository,
		fediverseService: fediverseService,
	}
}

// PostEcho 创建新的Echo
func (echoService *EchoService) PostEcho(userid uint, newEcho *model.Echo) error {
	err := echoService.txManager.Run(func(ctx context.Context) error {
		newEcho.UserID = userid

		user, err := echoService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return err
		}

		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		// 检查Extension内容
		if newEcho.Extension != "" && newEcho.ExtensionType != "" {
			switch newEcho.ExtensionType {
			case model.Extension_MUSIC:
				// 处理音乐链接 (暂无)
			case model.Extension_VIDEO:
				// 处理视频链接 (暂无)
			case model.Extension_GITHUBPROJ:
				// 处理GitHub项目的链接
				newEcho.Extension = httpUtil.TrimURL(newEcho.Extension)
			case model.Extension_WEBSITE:
				// 处理网站链接 (暂无)
			}
		} else {
			newEcho.Extension = ""
			newEcho.ExtensionType = ""
		}

		newEcho.Username = user.Username

		for i := range newEcho.Images {
			if newEcho.Images[i].ImageURL == "" {
				newEcho.Images[i].ImageSource = ""
			}
		}

		if newEcho.Content == "" && len(newEcho.Images) == 0 && (newEcho.Extension == "" || newEcho.ExtensionType == "") {
			return errors.New(commonModel.ECHO_CAN_NOT_BE_EMPTY)
		}

		// 处理临时文件表，防止被当作孤儿文件删除
		for i := range newEcho.Images {
			// 只有S3图片且有ObjectKey的才处理
			if newEcho.Images[i].ImageSource == model.ImageSourceS3 && newEcho.Images[i].ObjectKey != "" {
				// 直接删除临时文件记录 (开启事务)
				echoService.txManager.Run(func(ctx context.Context) error {
					return echoService.commonRepository.DeleteTempFileByObjectKey(ctx, newEcho.Images[i].ObjectKey)
				})
			}
		}

		return echoService.echoRepository.CreateEcho(ctx, newEcho)
	})

	if err != nil {
		return err
	}

	// 事务提交成功后再推送，确保已拿到持久化 ID
	savedEcho, fetchErr := echoService.echoRepository.GetEchosById(newEcho.ID)
	if fetchErr != nil {
		return fetchErr
	}
	if savedEcho != nil {
		if pushErr := echoService.fediverseService.PushEchoToFediverse(userid, *savedEcho); pushErr != nil {
			// 推送失败不影响发布
			fmt.Println("Error pushing Echo to Fediverse:", pushErr)
		}
	}

	return nil
}

// GetEchosByPage 获取Echo列表，支持分页
func (echoService *EchoService) GetEchosByPage(userid uint, pageQueryDto commonModel.PageQueryDto) (commonModel.PageQueryResult[[]model.Echo], error) {
	// 参数校验
	if pageQueryDto.Page < 1 {
		pageQueryDto.Page = 1
	}
	if pageQueryDto.PageSize < 1 || pageQueryDto.PageSize > 100 {
		pageQueryDto.PageSize = 10
	}

	//管理员登陆则支持查看隐私数据，否则不允许
	showPrivate := false
	if userid == authModel.NO_USER_LOGINED {
		showPrivate = false
	} else {
		user, err := echoService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return commonModel.PageQueryResult[[]model.Echo]{}, err
		}
		if user.IsAdmin {
			showPrivate = true
		} else {
			showPrivate = false
		}
	}

	echosByPage, total := echoService.echoRepository.GetEchosByPage(pageQueryDto.Page, pageQueryDto.PageSize, pageQueryDto.Search, showPrivate)
	result := commonModel.PageQueryResult[[]model.Echo]{
		Items: echosByPage,
		Total: total,
	}

	// 处理echosByPage中的图片URL
	for i := range result.Items {
		echoService.commonService.RefreshEchoImageURL(&result.Items[i])
	}

	// 返回结果
	return result, nil
}

// DeleteEchoById 删除指定ID的Echo
func (echoService *EchoService) DeleteEchoById(userid, id uint) error {
	return echoService.txManager.Run(func(ctx context.Context) error {
		user, err := echoService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		// 检查该Echo是否存在图片
		echo, err := echoService.echoRepository.GetEchosById(id)
		if err != nil {
			return err
		}
		if echo == nil {
			return errors.New(commonModel.ECHO_NOT_FOUND)
		}

		// 删除Echo中的图片
		if len(echo.Images) > 0 {
			for _, img := range echo.Images {
				if err := echoService.commonService.DirectDeleteImage(img.ImageURL, img.ImageSource, img.ObjectKey); err != nil {
					return err
				}
			}
		}

		return echoService.echoRepository.DeleteEchoById(ctx, id)
	})

}

// GetTodayEchos 获取今天的Echo列表
func (echoService *EchoService) GetTodayEchos(userid uint) ([]model.Echo, error) {
	//管理员登陆则支持查看隐私数据，否则不允许
	showPrivate := false
	if userid == authModel.NO_USER_LOGINED {
		showPrivate = false
	} else {
		user, err := echoService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return nil, err
		}
		if user.IsAdmin {
			showPrivate = true
		} else {
			showPrivate = false
		}
	}

	// 获取当日发布的Echos
	todayEchos := echoService.echoRepository.GetTodayEchos(showPrivate)

	// 处理todayEchos中的图片URL
	for i := range todayEchos {
		echoService.commonService.RefreshEchoImageURL(&todayEchos[i])
	}

	return todayEchos, nil
}

// UpdateEcho 更新指定ID的Echo
func (echoService *EchoService) UpdateEcho(userid uint, echo *model.Echo) error {
	return echoService.txManager.Run(func(ctx context.Context) error {
		user, err := echoService.commonService.CommonGetUserByUserId(userid)
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}

		// 检查Extension内容
		if echo.Extension != "" && echo.ExtensionType != "" {
			switch echo.ExtensionType {
			case model.Extension_MUSIC:
				// 处理音乐链接 (暂无)
			case model.Extension_VIDEO:
				// 处理视频链接 (暂无)
			case model.Extension_GITHUBPROJ:
				echo.Extension = httpUtil.TrimURL(echo.Extension)
			case model.Extension_WEBSITE:
				// 处理网站链接 (暂无)
			}
		} else {
			echo.Extension = ""
			echo.ExtensionType = ""
		}

		// 处理无效图片
		for i := range echo.Images {
			if echo.Images[i].ImageURL == "" {
				echo.Images[i].ImageSource = ""
				echo.Images[i].ImageURL = ""
			}
			// 确保外键正确设置
			echo.Images[i].MessageID = echo.ID
		}

		// 检查是否为空
		if echo.Content == "" && len(echo.Images) == 0 && (echo.Extension == "" || echo.ExtensionType == "") {
			return errors.New(commonModel.ECHO_CAN_NOT_BE_EMPTY)
		}

		return echoService.echoRepository.UpdateEcho(ctx, echo)
	})

}

// LikeEcho 点赞指定ID的Echo
func (echoService *EchoService) LikeEcho(id uint) error {
	return echoService.txManager.Run(func(ctx context.Context) error {
		return echoService.echoRepository.LikeEcho(ctx, id)
	})

}

// GetEchoById 获取指定 ID 的 Echo
func (echoService *EchoService) GetEchoById(userId, id uint) (*model.Echo, error) {
	var echo *model.Echo

	echo, err := echoService.echoRepository.GetEchosById(id)
	if err != nil {
		return nil, err
	}

	// 如果不存在Echo，则返回错误
	if echo == nil {
		return nil, errors.New(commonModel.ECHO_NOT_FOUND)
	}

	// 如果没有登录用户，则不允许获取私密Echo
	if userId == authModel.NO_USER_LOGINED {
		// 如果Echo是私密的，则不允许获取
		if echo.Private {
			// 不允许通过ID获取私密Echo
			return nil, errors.New(commonModel.NO_PERMISSION_DENIED)
		}
	} else {
		// 如果用户已经登录,获取当前登录用户
		user, err := echoService.commonService.CommonGetUserByUserId(userId)
		if err != nil {
			return nil, err
		}

		if echo.Private {
			if !user.IsAdmin {
				return nil, errors.New(commonModel.NO_PERMISSION_DENIED)
			}

			return echo, nil
		}
	}

	// 刷新图片URL
	echoService.commonService.RefreshEchoImageURL(echo)

	// 返回Echo
	return echo, nil
}

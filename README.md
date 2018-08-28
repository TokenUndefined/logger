# logger
基于GO语言的服务端日志模块，接口简单明了，易嵌入到工程中去。
# 特性
1.支持按日切分；
2.支持按大小切分，默认大小为100M；
3.支持即按日又按大小切分log；
4.支持控制台不同级别日志显示不同颜色；
# 获取
go get github.com/TokenUndefined/logger
# 使用
import github.com/TokenUndefined/logger
logger.Debug("hai")
logger.Info("hai info msg")
logger.Warn("hai Warn")
logger.Error("hai Error msg")
logger.Fatal("hai Fatal msg")

logger.Debugf("I'm %s log! ","debug")
logger.Infof("I'm %s log!","info")
logger.Warnf("I'm %s log!","warn")
logger.Errorf("I'm %s log!","error")

logger.Debugln("I'm","debug","log!")
logger.Infoln("I'm","info","log!")
logger.Warnln("I'm","warn","log!")
logger.Errorln("I'm","error","log!")

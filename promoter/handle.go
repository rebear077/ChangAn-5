package promote

import (
	"fmt"
	"sync"
	"time"

	server "github.com/rebear077/changan/backend"
	receive "github.com/rebear077/changan/connApi"
	"github.com/rebear077/changan/errorhandle"
	logloader "github.com/rebear077/changan/logs"
	uptoChain "github.com/rebear077/changan/tochain"
	"github.com/sirupsen/logrus"
)

var logs = logloader.NewLog()

type Promoter struct {
	server        *server.Server
	DataApi       *receive.FrontEnd
	encryptedPool *Pools
	loader        *logloader.Loader
}

func NewPromoter() *Promoter {
	ser := server.NewServer()
	api := receive.NewFrontEnd()
	pool := NewPools()
	lder := logloader.NewLoader()
	// chainld := chainloader.NewChainInfo()
	return &Promoter{
		server:        ser,
		DataApi:       api,
		encryptedPool: pool,
		loader:        lder,
	}
}

func (p *Promoter) Start() {
	go p.loader.Start()
	go p.server.StartMonitor()
	logs.Infoln("开始运行")
	for {
		if p.server.VerifyChainstatus() {
			p.InvoiceInfoHandler()
			p.SupplierFinancingApplicationInfoHandler()
			p.HistoricalInfoHandler()
			p.PushPaymentAccountsInfoHandler()
			p.PoolInfoHandler()
		} else {
			time.Sleep(10 * time.Second)
		}
	}

}

func (p *Promoter) InvoiceInfoHandler() {
	if len(p.DataApi.InvoicePool) != 0 {
		// logrus.Infoln(len(p.DataApi.InvoicePool))
		logs.Infoln(len(p.DataApi.InvoicePool))
		// logrus.Infoln("开始同步发票信息")
		logs.Infoln("开始同步发票信息")
		var wg sync.WaitGroup
		invoices := make([]*receive.InvoiceInformation, 0)
		p.DataApi.Invoicemutex.Lock()
		invoices = append(invoices, p.DataApi.InvoicePool...)
		p.DataApi.InvoicePool = nil
		p.DataApi.Invoicemutex.Unlock()
		mapping := server.EncodeInvoiceInformation(invoices)
		// logrus.Infoln(len(mapping))
		logs.Infoln(len(mapping))
		for index := range mapping {
			for header, info := range mapping[index] {
				wg.Add(1)
				tempheader := header
				tempinfo := info
				go func(tempheader string, tempinfo string) {
					p.packInfo(tempheader, tempinfo, "fast", "invoice")
					wg.Done()
				}(tempheader, tempinfo)
			}
		}
		wg.Wait()
		messages := p.encryptedPool.QueryMessages("invoice", "fast")
		for _, message := range messages {
			temp, _ := message.(packedMessage)
			err := p.server.IssueInvoiceInformation(temp.header, temp.cipher, temp.encryptionKey, temp.signed)
			if err != nil {
				// logrus.Errorln("发票信息上链失败:", temp.header, "失败信息为:", err)
				logs.Errorln("发票信息上链失败:", temp.header, "失败信息为:", err)
			}
		}
		for {
			errNum := errorhandle.ERRDealer.GetInvoiceInfoPoolLength()
			success := uptoChain.QueryInvoiceSuccessCounter()
			if errNum+success == len(messages) {
				// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				uptoChain.ResetInvoiceSuccessCounter()
				break
			}
		}
	}
}

func (p *Promoter) HistoricalInfoHandler() {
	if len(p.DataApi.TransactionHistoryPool) != 0 {
		// logrus.Infoln("开始历史交易信息")
		logs.Infoln("开始历史交易信息")
		var wg sync.WaitGroup
		hisinfos := make([]*receive.TransactionHistory, 0)
		p.DataApi.TransactionHistorymutex.Lock()
		hisinfos = append(hisinfos, p.DataApi.TransactionHistoryPool...)
		p.DataApi.TransactionHistoryPool = nil
		p.DataApi.TransactionHistorymutex.Unlock()
		mapping := server.EncodeTransactionHistory(hisinfos)
		for index := range mapping {
			for header, info := range mapping[index] {
				tempheader := header
				tempinfo := info
				wg.Add(1)
				go func(tempheader string, tempinfo string) {
					usedvalue, settlevalue, ordervalue, receivablevalue := server.HistoricalInformationSlice(tempheader, tempinfo, 5)
					p.packInfos(tempheader, usedvalue, "fast", "historicalUsed")
					p.packInfos(tempheader, settlevalue, "fast", "historicalSettle")
					p.packInfos(tempheader, ordervalue, "fast", "historicalOrder")
					p.packInfos(tempheader, receivablevalue, "fast", "historicalReceivable")
					wg.Done()
				}(tempheader, tempinfo)
			}
		}
		wg.Wait()
		hisUsedMessage := p.encryptedPool.QueryMessages("historicalUsed", "fast")
		hisSettleMessage := p.encryptedPool.QueryMessages("historicalSettle", "fast")
		hisOrderMessage := p.encryptedPool.QueryMessages("historicalOrder", "fast")
		hisReceivableMessage := p.encryptedPool.QueryMessages("historicalReceivable", "fast")
		for _, message := range hisUsedMessage {
			tempUsed, _ := message.(packedMessage)
			err := p.server.IssueHistoricalUsedInformation(tempUsed.header, tempUsed.cipher, tempUsed.encryptionKey, tempUsed.signed)
			if err != nil {
				// logrus.Errorln("信息上链失败:", tempUsed.header, "失败信息为:", err)
				logs.Errorln("信息上链失败:", tempUsed.header, "失败信息为:", err)
			}
		}
		for _, message := range hisSettleMessage {
			tempSettle, _ := message.(packedMessage)
			err := p.server.IssueHistoricalSettleInformation(tempSettle.header, tempSettle.cipher, tempSettle.encryptionKey, tempSettle.signed)
			if err != nil {
				// logrus.Errorln("信息上链失败:", tempSettle.header, "失败信息为:", err)
				logs.Errorln("信息上链失败:", tempSettle.header, "失败信息为:", err)
			}
		}
		for _, message := range hisOrderMessage {
			tempOrder, _ := message.(packedMessage)
			err := p.server.IssueHistoricalOrderInformation(tempOrder.header, tempOrder.cipher, tempOrder.encryptionKey, tempOrder.signed)
			if err != nil {
				// logrus.Errorln("信息上链失败:", tempOrder.header, "失败信息为:", err)
				logs.Errorln("信息上链失败:", tempOrder.header, "失败信息为:", err)
			}
		}
		for _, message := range hisReceivableMessage {
			tempReceivable, _ := message.(packedMessage)
			err := p.server.IssueHistoricalReceivableInformation(tempReceivable.header, tempReceivable.cipher, tempReceivable.encryptionKey, tempReceivable.signed)
			if err != nil {
				// logrus.Errorln("信息上链失败:", tempReceivable.header, "失败信息为:", err)
				logs.Errorln("信息上链失败:", tempReceivable.header, "失败信息为:", err)
			}
		}
		wg.Add(4)
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetHistoricalUsedInfoPoolLength()
				success := uptoChain.QueryHistoricalUsedCounter()
				if errNum+success == len(hisUsedMessage) {
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisUsedMessage), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisUsedMessage), success, errNum)
					uptoChain.ResetHistoricalUsedCounter()
					break
				}
			}
			wg.Done()
		}()
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetHistoricalSettleInfoPoolLength()
				success := uptoChain.QueryHistoricalSettleCounter()
				if errNum+success == len(hisSettleMessage) {
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisSettleMessage), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisSettleMessage), success, errNum)
					uptoChain.ResetHistoricalSettleCounter()
					break
				}
			}
			wg.Done()
		}()
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetHistoricalOrderInfoPoolLength()
				success := uptoChain.QueryHistoricalOrderCounter()
				if errNum+success == len(hisOrderMessage) {
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisOrderMessage), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisOrderMessage), success, errNum)
					uptoChain.ResetHistoricalOrderCounter()
					break
				}
			}
			wg.Done()
		}()
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetHistoricalReceivableInfoPoolLength()
				success := uptoChain.QueryHistoricalReceivableCounter()
				if errNum+success == len(hisReceivableMessage) {
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisReceivableMessage), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(hisReceivableMessage), success, errNum)
					uptoChain.ResetHistoricalReceivableCounter()
					break
				}
			}
			wg.Done()
		}()
		wg.Wait()

	}
}
func (p *Promoter) PoolInfoHandler() {
	if len(p.DataApi.EnterpoolDataPool) != 0 {
		// logrus.Infoln("开始入池数据信息")
		logs.Infoln("开始入池数据信息")
		var wg sync.WaitGroup
		poolinfos := make([]*receive.EnterpoolData, 0)
		p.DataApi.EnterpoolDatamutex.Lock()
		poolinfos = append(poolinfos, p.DataApi.EnterpoolDataPool...)
		p.DataApi.EnterpoolDataPool = nil
		p.DataApi.EnterpoolDatamutex.Unlock()
		mapping := server.EncodeEnterpoolData(poolinfos)
		for index := range mapping {
			for header, info := range mapping[index] {
				tempheader := header
				tempinfo := info
				wg.Add(1)
				go func(tempheader string, tempinfo string) {
					var wwg sync.WaitGroup
					planvalue, providerusedvalue := server.PoolInformationSlice(tempheader, tempinfo, 5)
					wwg.Add(2)
					go func(tempheader string, planvalue []string) {
						p.packInfos(tempheader, planvalue, "fast", "poolPlan")
						wwg.Done()
					}(tempheader, planvalue)
					go func(tempheader string, providerusedvalue []string) {
						p.packInfos(tempheader, providerusedvalue, "fast", "poolUsed")
						wwg.Done()
					}(tempheader, providerusedvalue)
					wwg.Wait()
					wg.Done()
				}(tempheader, tempinfo)
			}
		}
		wg.Wait()
		planMessages := p.encryptedPool.QueryMessages("poolPlan", "fast")
		usedMessages := p.encryptedPool.QueryMessages("poolUsed", "fast")
		for _, message := range planMessages {
			tempPlan, _ := message.(packedMessage)
			err := p.server.IssuePoolPlanInformation(tempPlan.header, tempPlan.cipher, tempPlan.encryptionKey, tempPlan.signed)
			if err != nil {
				// logrus.Errorln("信息上链失败:", tempPlan.header, "失败信息为:", err)
				logs.Errorln("信息上链失败:", tempPlan.header, "失败信息为:", err)
			}
		}
		for _, message := range usedMessages {
			tempUsed, _ := message.(packedMessage)
			err := p.server.IssuePoolUsedInformation(tempUsed.header, tempUsed.cipher, tempUsed.encryptionKey, tempUsed.signed)
			if err != nil {
				// logrus.Errorln("信息上链失败:", tempUsed.header, "失败信息为:", err)
				logs.Errorln("信息上链失败:", tempUsed.header, "失败信息为:", err)
			}

		}
		wg.Add(2)
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetPoolPlanInfoPoolLength()
				success := uptoChain.QueryPoolPlanCounter()
				if errNum+success == len(planMessages) {
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(planMessages), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(planMessages), success, errNum)
					uptoChain.ResetPoolPlanCounter()
					break
				}
			}
			wg.Done()
		}()
		go func() {
			for {
				errNum := errorhandle.ERRDealer.GetPoolUsedInfoPoolLength()
				success := uptoChain.QueryPoolUsedCounter()
				if errNum+success == len(usedMessages) {
					// logrus.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(usedMessages), success, errNum)
					logs.Infof("同步完成，共计%d条数据，成功%d,失败%d", len(usedMessages), success, errNum)
					uptoChain.ResetPoolUsedCounter()
					break
				}
			}
			wg.Done()
		}()
		wg.Wait()
		// logrus.Println("退出")
		logs.Println("退出")
	}
}

func (p *Promoter) SupplierFinancingApplicationInfoHandler() {
	if len(p.DataApi.FinancingIntentionPool) != 0 {
		// logrus.Infoln("开始同步融资意向请求信息")
		logs.Infoln("开始同步融资意向请求信息")
		var wg sync.WaitGroup
		finintens := make([]*receive.FinancingIntention, 0)
		p.DataApi.FinancingIntentionmutex.Lock()
		finintens = append(finintens, p.DataApi.FinancingIntentionPool...)
		p.DataApi.FinancingIntentionPool = nil
		p.DataApi.FinancingIntentionmutex.Unlock()
		mapping := server.EncodeFinancingIntention(finintens)
		for index := range mapping {
			for header, info := range mapping[index] {
				wg.Add(1)
				tempheader := header
				tempinfo := info
				go func(tempheader string, tempinfo string) {
					p.packInfo(tempheader, tempinfo, "fast", "application")
					wg.Done()
				}(tempheader, tempinfo)
			}
		}
		wg.Wait()
		messages := p.encryptedPool.QueryMessages("application", "fast")
		for _, message := range messages {
			temp, _ := message.(packedMessage)
			err := p.server.IssueSupplierFinancingApplication(temp.header, temp.cipher, temp.encryptionKey, temp.signed)
			if err != nil {
				// logrus.Errorln("融资意向请求上链失败,", "失败信息为:", err)
				logs.Errorln("融资意向请求上链失败,", "失败信息为:", err)
			}
		}
		for {
			errNum := errorhandle.ERRDealer.GetSupplierFinancingApplicationPoolLength()
			success := uptoChain.QuerySupplierSuccessCounter()
			if errNum+success == len(messages) {
				// logrus.Infof("同步融资意向完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				logs.Infof("同步融资意向完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				uptoChain.ResetSupplierSuccessCounter()
				break
			}
		}
	}
}

func (p *Promoter) PushPaymentAccountsInfoHandler() {
	if len(p.DataApi.CollectionAccountPool) != 0 {
		logrus.Infoln("开始同步回款信息")
		logs.Infoln("开始同步回款信息")
		var wg sync.WaitGroup
		payinfos := make([]*receive.CollectionAccount, 0)
		p.DataApi.CollectionAccountmutex.Lock()
		payinfos = append(payinfos, p.DataApi.CollectionAccountPool...)
		p.DataApi.CollectionAccountPool = nil
		p.DataApi.CollectionAccountmutex.Unlock()
		mapping := server.EncodeCollectionAccount(payinfos)
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
		for index := range mapping {
			for header, info := range mapping[index] {
				wg.Add(1)
				tempheader := header
				tempinfo := info
				go func(tempheader string, tempinfo string) {
					p.packInfo(tempheader, tempinfo, "fast", "payment")
					wg.Done()
				}(tempheader, tempinfo)
			}
		}
		wg.Wait()
		fmt.Println("dddddddddddddddddddddddddddddddd")
		messages := p.encryptedPool.QueryMessages("payment", "fast")
		fmt.Println(len(messages))
		for index, message := range messages {
			fmt.Println(index)
			temp, ok := message.(packedMessage)
			if !ok {
				fmt.Println("errorerror")
			}

			err := p.server.IssuePushPaymentAccount(temp.header, temp.cipher, temp.encryptionKey, temp.signed)
			if err != nil {
				// logrus.Errorln("回款信息上链失败,", "失败信息为:", err)
				logs.Errorln("回款信息上链失败,", "失败信息为:", err)
			}
		}
		fmt.Println("jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj")
		for {
			errNum := errorhandle.ERRDealer.GetPushPaymentAccountPoolLength()
			success := uptoChain.QueryPaymentAccountsCounter()
			if errNum+success == len(messages) {
				// logrus.Infof("回款信息同步完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				logs.Infof("回款信息同步完成，共计%d条数据，成功%d,失败%d", len(messages), success, errNum)
				uptoChain.ResetPaymentAccountsCounter()
				break
			}
		}
	}
}
func (p *Promoter) packInfo(header string, info string, poolType string, method string) {
	cipher, encryptionKey, signed, err := p.server.DataEncryption([]byte(info))
	if err != nil {
		// logrus.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
		logs.Fatalln("数据加密失败,此条数据信息为:", header, info, "失败信息为:", err)
	}
	temp := packedMessage{}
	temp.cipher = cipher
	temp.encryptionKey = encryptionKey
	temp.signed = signed
	temp.header = header
	p.encryptedPool.Insert(temp, method, poolType)
}
func (p *Promoter) packInfos(header string, infos []string, poolType string, method string) {
	var wg sync.WaitGroup
	for _, info := range infos {
		tempinfo := info
		wg.Add(1)
		go func(header string, tempinfo string) {
			cipher, encryptionKey, signed, err := p.server.DataEncryption([]byte(tempinfo))
			if err != nil {
				// logrus.Fatalln("数据加密失败,此条数据信息为:", header, tempinfo, "失败信息为:", err)
				logs.Fatalln("数据加密失败,此条数据信息为:", header, tempinfo, "失败信息为:", err)
			}
			temp := packedMessage{}
			temp.cipher = cipher
			temp.encryptionKey = encryptionKey
			temp.signed = signed
			temp.header = header
			p.encryptedPool.Insert(temp, method, poolType)
			wg.Done()
		}(header, tempinfo)

	}
	wg.Wait()
}

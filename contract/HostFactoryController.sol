pragma solidity ^0.4.25;
pragma experimental ABIEncoderV2;
import "./Ownable.sol";
import "./SupplierFinancingApplication.sol";
import "./InvoiceInformationStorage.sol";
import "./PushPaymentAccounts.sol";
import "./PublicKeyStorage.sol";
import "./HistoryOrderInfo.sol";
import "./HistoryReceivableInfo.sol";
import "./HistoryUsedInfo.sol";
import "./HistorySettleInfo.sol";
import "./PoolPlanInfo.sol";
import "./PoolUsedInfo.sol";
contract HostFactoryController is Ownable{

    SupplierFinancingApplication private supplierFinancingApplication;
    InvoiceInformationStorage private invoiceInformationStorage;
    PushPaymentAccounts private pushPaymentAccounts;
    PublicKeyStorage private publicKeyStorage;
    HistoryOrderInfo private historyOrderInfo;
    HistoryReceivableInfo private historyReceivableInfo;
    HistoryUsedInfo private historyUsedInfo;
    HistorySettleInfo private historySettleInfo;
    PoolPlanInfo private poolPlanInfo;
    PoolUsedInfo private poolUsedInfo;
    
    constructor() public {
        supplierFinancingApplication=new SupplierFinancingApplication();
        invoiceInformationStorage=new InvoiceInformationStorage();
        pushPaymentAccounts=new PushPaymentAccounts();
        publicKeyStorage=new PublicKeyStorage();
        historyOrderInfo=new HistoryOrderInfo();
        historyReceivableInfo=new HistoryReceivableInfo();
        historyUsedInfo =new HistoryUsedInfo();
        historySettleInfo =new HistorySettleInfo();
        poolPlanInfo=new PoolPlanInfo();
        poolUsedInfo=new PoolUsedInfo();
        
    }
//*********************************************************************************************
//公钥   
    function issuePublicKeyStorage(string _id,string _role,string _key)external onlyOwner returns(int256){
        int256 count = publicKeyStorage.insert(_id, _role,_key);
        return count;
    }
    function queryPublicKey(string _id) public returns(string){
        string memory res= publicKeyStorage.getDetail(_id);
        return res;
    }  
//***************************************************************************************************    
//融资意向    
    function issueSupplierFinancingApplication(string _id, string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = supplierFinancingApplication.insert(_id, _data,_key,_hash);
        
        return count;
    }
    //返回列表格式
    function querySupplierFinancingApplicationInList(string _id)external returns(string){
        string memory res=supplierFinancingApplication.getDetail(_id);
        return res;
    }
    //返回json格式
    function querySupplierFinancingApplicationInJson(string _id)external returns(string){
        string memory res=supplierFinancingApplication.getDetailInJson(_id);
        return res;
    }
//*****************************************************************************************************
//发票信息  
    function issueInvoiceInformationStorage(string _id, string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = invoiceInformationStorage.insert(_id, _data,_key,_hash);
      
        return count;
    }
    function queryInvoiceInformationInList(string _id)public returns(string){
        string memory res=invoiceInformationStorage.getDetailInList(_id);
        return res;
    }
    function queryInvoiceInformationInJson(string _id)public returns(string){
        string memory res=invoiceInformationStorage.getDetailInJson(_id);
        return res;
    }
//*******************************************************************************************************    
//历史交易信息之入库信息
    function issueHistoricalUsedInformation(string _id, string _time,string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = historyUsedInfo.insert(_id,_time,_data,_key,_hash);
        return count;
    }
    function queryHIstoricalUsedInList(string _id)public returns(string){
        string memory res=historyUsedInfo.getDetailInList(_id);
        return res;
    }
    function queryHIstoricalUsedInJson(string _id)public returns(string){
        string memory res=historyUsedInfo.getDetailInJson(_id);
        return res;
    }
//*******************************************************************************************************    
//历史交易信息之订单信息
    function issueHistoricalOrderInformation(string _id, string _time,string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = historyOrderInfo.insert(_id,_time,_data,_key,_hash);
        return count;
    }
    function queryHIstoricalOrderInList(string _id)public returns(string){
        string memory res=historyOrderInfo.getDetailInList(_id);
        return res;
    }
    function queryHIstoricalOrderInJson(string _id)public returns(string){
        string memory res=historyOrderInfo.getDetailInJson(_id);
        return res;
    }
//*******************************************************************************************************    
//历史交易信息之结算信息
    function issueHistoricalSettleInformation(string _id, string _time,string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = historySettleInfo.insert(_id,_time,_data,_key,_hash);
        return count;
    }
    function queryHIstoricalSettleInList(string _id)public returns(string){
        string memory res=historySettleInfo.getDetailInList(_id);
        return res;
    }
    function queryHIstoricalSettleInJson(string _id)public returns(string){
        string memory res=historySettleInfo.getDetailInJson(_id);
        return res;
    }
//*******************************************************************************************************    
//历史交易信息之应收账款信息
    function issueHistoricalReceivableInformation(string _id, string _time,string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = historyReceivableInfo.insert(_id,_time,_data,_key,_hash);
        return count;
    }
    function queryHIstoricalReceivableInList(string _id)public returns(string){
        string memory res=historyReceivableInfo.getDetailInList(_id);
        return res;
    }
    function queryHIstoricalReceivableInJson(string _id)public returns(string){
        string memory res=historyReceivableInfo.getDetailInJson(_id);
        return res;
    }
//********************************************************************************************************
//查询回款信息
    function issuePushPaymentAccounts(string _id, string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = pushPaymentAccounts.insert(_id, _data,_key,_hash);
        return count;
    }
    function queryPushPaymentAccountsInList(string _id)public returns(string){
        string memory res=pushPaymentAccounts.getDetail(_id);
        return res;
    }
    function queryPushPaymentAccountsInJson(string _id)public returns(string){
        string memory res=pushPaymentAccounts.getDetailInJson(_id);
        return res;
    }
    
//*******************************************************************************************************    
//查询入池数据之供应商生产计划信息
    function issuePoolPlanInformation(string _id, string _time,string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = poolPlanInfo.insert(_id,_time,_data,_key,_hash);
        return count;
    }
    function queryPoolPlanInfoInList(string _id)public returns(string){
        string memory res=poolPlanInfo.getDetailInList(_id);
        return res;
    }
    function queryPoolPlanInfoInJson(string _id)public returns(string){
        string memory res=poolPlanInfo.getDetailInJson(_id);
        return res;
    }
//*******************************************************************************************************    
//查询入池数据之供应商生产入库信息
    function issuePoolUsedInformation(string _id, string _time,string _data,string _key,string _hash) external onlyOwner returns(int256){
        int256 count = poolUsedInfo.insert(_id,_time,_data,_key,_hash);
        return count;
    }
    function queryPoolUsedInfoInList(string _id)public returns(string){
        string memory res=poolUsedInfo.getDetailInList(_id);
        return res;
    }
    function queryPoolUsedInfoInJson(string _id)public returns(string){
        string memory res=poolUsedInfo.getDetailInJson(_id);
        return res;
    }
}
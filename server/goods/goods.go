package goods

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"sync"
	"zdshop/models"
	pb "zdshop/proto-go"
)

type goodsServer struct {
	pb.UnimplementedGoodsServer
}

func NewGoodsServer() *goodsServer {
	return &goodsServer{}
}

type goodsInfo struct {
	goods models.Goods
	parameters []*models.GoodsParameter
	banner []*models.GoodsBanner
}

func newGoodsInfo() *goodsInfo {
	var goods models.Goods
	var parameters []*models.GoodsParameter
	var banner []*models.GoodsBanner

	info := &goodsInfo{
		goods:      goods,
		parameters: parameters,
		banner:     banner,
	}
	return info
}

func (g *goodsServer) GetOneGoods(ctx context.Context, req *pb.GetOneGoodsRequest) (*pb.GetOneGoodsReply, error) {
	fmt.Println("grpc get one goods")
	id := req.Gid

	sy := sync.WaitGroup{}
	sy.Add(1)
	goodsInfo := newGoodsInfo()
	go func() {
		o := orm.NewOrm()
		_ := o.QueryTable("Goods").Filter("Id", id).One(&goodsInfo.goods, "Id", "Name", "GoodsBrief", "Explain", "SalesValue", "category_id")
		_, _ = o.QueryTable("GoodsParameter").Filter("Goods__Id", id).All(&goodsInfo.parameters, "Id", "Parameter", "GoodsTruePrice", "GoodsNowPrice", "IsDefault", "GoodsNumber")
		_, _ = o.QueryTable("GoodsBanner").Filter("Goods__Id", id).OrderBy("Id").All(&goodsInfo.banner, "Id", "GoodsUrl", "IsDefault")
		sy.Done()
	}()

	sy.Wait()
	var pbGoods pb.GoodsInfo
	pbGoods.Id = int64(goodsInfo.goods.Id)
	pbGoods.Name = goodsInfo.goods.Name
	pbGoods.GoodsBrief = goodsInfo.goods.GoodsBrief
	pbGoods.Explain = goodsInfo.goods.Explain
	pbGoods.SalesValue = strconv.Itoa(goodsInfo.goods.SalesValue)
	pbGoods.CategoryId = strconv.Itoa(goodsInfo.goods.Category.Id)

	return &pb.GetOneGoodsReply{
		Goods:      &pbGoods,
		Parameters: nil,
		Banner:     nil,
	}, nil
}

func test() {

	//o := orm.NewOrm()
	//var goods models.Goods
	//_ := o.QueryTable("Goods").Filter("Id", id).One(&goods, "Id", "Name", "GoodsBrief", "Explain", "SalesValue", "category_id")
	//if err != nil {
	//	return nil, status.Error(codes.NotFound, "商品查询失败")
	//}
	//var parameters []*models.GoodsParameter
	//_, _ = o.QueryTable("GoodsParameter").Filter("Goods__Id", id).All(&parameters, "Id", "Parameter", "GoodsTruePrice", "GoodsNowPrice", "IsDefault", "GoodsNumber")
	//if err != nil {
	//	return nil, status.Error(codes.NotFound, "商品参数查询失败")
	//}
	//var banner []*models.GoodsBanner
	//_, _ = o.QueryTable("GoodsBanner").Filter("Goods__Id", id).OrderBy("Id").All(&banner, "Id", "GoodsUrl", "IsDefault")
	//if err != nil {
	//	return nil, status.Error(codes.NotFound, "商品参数查询失败")
	//}
}

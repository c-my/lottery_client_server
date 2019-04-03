window.onload=function(){

    //获取图片列表，之后后台在这给数据就行
        function pictureList(){
            var pictureList=[];
            for(var i=1;i<=12;i++){
              pictureList.push("img/"+i+".jpeg");//把得到的图片放到数组中
            }
              return pictureList;//返回图片数组
        };

        var list=pictureList();//图片数组
        var divList=[];//div数组
        for(var i in list){
       //添加图片,
        var newDiv=document.createElement("div");
        var newImg=document.createElement("img");
        //var newSDiv=document.createElement("div");
        var user=document.getElementById("swing");
        newDiv.className="userW";
        newImg.className="pro";
       // newSDiv.className="userS";
        newDiv.appendChild(newImg);
       // newDiv.appendChild(newSDiv);
        newImg.src=list[i];
        user.appendChild(newDiv);
        divList.push(newImg);//获取div数组

    }
        console.log(divList);
        //var flag=false;
     $.fn.draw = function(options){
        var $_this = $(this);
        var $dp = options.prod;
        var $i = 0; //index
        var $r = 0; //round
        var $s = 150; //spead
        //快速闪动，会走两遍
        function dr(){
            //drawround
          if($_this.res!='' && $_this.res!=undefined && $i==0){
                  dr2();
          }else{
            $dp.find(".pro").removeClass("select").eq($i).addClass("select");
            $i = $i >= 11 ? 0 : $i+1;//用户数组的数量
            setTimeout(dr,$s);
          }
        }
          function dr2(){
         $dp.find(".pro").removeClass("select").eq($i).addClass("select");
          $i = $i >= 11 ? 0 : $i+1;
          $s = $s+30;//让他变慢
          $rand=Math.ceil(Math.random()*12)
          if( $r < $_this.res + $rand ){
            $r++;
          }else{
            $i = 0;
            $r = 0;
            $s = 100;
            setTimeout(result,1000);
            return;
          }
              setTimeout(dr2,$s);
          }

        function getRes(){
            $_this.res = '';
              //赋结果值
            setTimeout(function(){$_this.res = 1},2000);
        }

        //监听点击事件
        function click(){
          $_this.bind("click",function(){
            getRes();
            console.log(i);
            dr();
            $_this.unbind("click");
          });
        }click();

        //输出获奖者
        function result(){
          alert("XX得到奖项");
          click();
        }
      }
      $("#btn").draw({
        prod:$("#swing")
    });
}

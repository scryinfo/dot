<template>
  <div>
    <el-container id="exportDot" style="margin-top: -12px;">
      <el-header id="eDotH" style="line-height: 30px;height: 30px;">Dot</el-header>
      <el-main id="eDotM">
        <div>
          <span>fileName :</span>
          <el-input type="text" v-model="dotFileName" style="width: 60%;margin-left: 2%;"></el-input>
        </div>
        <el-button
          id="expoertD"
          @click="ExportDot"
          style="margin-top: 20px;margin-left: 26%;"
        >Export Dot</el-button>
      </el-main>
    </el-container>
    <el-container id="exportConf">
      <el-header id="eConfH" style="height: 30px;line-height: 30px">Config</el-header>
      <el-main id="eConfM">
        <div>
          <span>fileName :</span>
          <el-input type="text" v-model="confFileName" style="width: 60%;margin-left: 2%;"></el-input>
        </div>
        <el-row style="margin-top: 20px">
          <template>
            <el-checkbox
              :indeterminate="isIndeterminate"
              v-model="checkAllC"
              @change="handleCheckAllChangeC"
            >全选</el-checkbox>
            <div style="margin: 15px 0;"></div>
            <el-checkbox-group v-model="checkedCitiesC" @change="handleCheckedCitiesChangeC">
              <el-checkbox v-for="city in cities" :label="city" :key="city">{{city}}</el-checkbox>
            </el-checkbox-group>
          </template>
        </el-row>
        <el-button
          id="findC"
          @click="ExportConf"
          style="margin-top: 20px;margin-left: 26%;"
        >ExportConfig</el-button>
      </el-main>
    </el-container>
  </div>
</template>

<script>
const cityOptions = ["json", "toml", "yaml"];
export default {
  data() {
    return {
      dotFileName: "",
      confFileName: "",
      checkAllC: false,
      checkedCitiesC: ["json"],
      cities: cityOptions,
      isIndeterminate: true
    };
  },
  methods: {
    handleCheckAllChangeC(val) {
      this.checkedCitiesC = val ? cityOptions : [];
      this.isIndeterminate = false;
    },
    handleCheckedCitiesChangeC(value) {
      let checkedCount = value.length;
      this.checkAllC = checkedCount === this.cities.length;
      this.isIndeterminate =
        checkedCount > 0 && checkedCount < this.cities.length;
    },
    ExportDot() {
      if (this.dotFileName == "") {
        alert("请输入文件名");
      } else {
        var dotfilename = this.dotFileName + ".json";
        var filename = [dotfilename];
        var { rpcExportDot } = require("../plugins/rpcInterface");
        rpcExportDot(this.$root.Dots, filename, response => {
          if (response.getError() == "") {
            alert(
              "导出文件" +
                filename +
                "成功，文件位置tools/config/data/server目录下"
            );
          } else {
            alert("导出文件" + filenames + "失败" + response.getError());
            console.log(response);
          }
          console.log(this.$root.Dots);
        });
      }
    },
    ExportConf() {
      if (this.confFileName == "") {
        alert("请输入文件名");
      } else {
        var confFileNames = [];
        //要生成的文件名
        for (var i = 0; i < this.checkedCitiesC.length; i++) {
          confFileNames.push(this.confFileName + "." + this.checkedCitiesC[i]);
        }
        //判断liveid
        var conf = this.$root.Configs; //config页面数据
        var resultDot = []; //处理掉空配置
        
        {
          var liveIds = [];
          for (var i = 0; i < conf.length; i++) {
            if (conf[i].lives.length == 0) {
              //实例数为０跳过
              continue;
            }
            for (var j = 0; j < conf[i].lives.length; j++) {
              if (conf[i].lives[j].liveId == "") {
                alert(conf[i].lives[j] + ":liveId is null");
                return false;
              } else {
                for (var z = 0; z < liveIds.length; z++) {
                  if (conf[i].lives[j].liveId == liveIds[z]) {
                    alert(conf[i].lives[j].liveId + "liveid重复．");
                    return false;
                  }
                }
                liveIds.push(conf[i].lives[j].liveId);
              }
            }
            resultDot.push(conf[i]);
          }
        }
        var result = {
          log: {
            file: "log.log",
            level: "debug"
          },
          dots: null
        };
        result.dots = resultDot;
        var { rpcExportConfig } = require("../plugins/rpcInterface");
        rpcExportConfig(result, confFileNames, response => {
          if (response.getError() == "") {
            alert(
              "导出文件" +
                confFileNames +
                "成功，文件位置tools/config/data/server目录下"
            );
          } else {
            alert("导出文件" + confFileNames + "失败" + response.getError());
          }
        });
      }
    }
  }
};
</script>
<style>
#eDotH,
#eConfH {
  text-align: left;
  background-color: #d3dce6;
}
</style>

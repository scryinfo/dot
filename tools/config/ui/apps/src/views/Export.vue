<template>
  <div >
    <el-container id="exportDot" style="margin-top: -12px;">
      <el-header id="eDotH" style="line-height: 30px;height: 30px;">Dot</el-header>
      <el-main id="eDotM">
        <div >
          <span>fileName : </span>
          <el-input type="text"  v-model="dotFileName" style="width: 60%;margin-left: 2%;"></el-input>
        </div>
        <el-button id="expoertD" @click="ExportDot" style="margin-top: 20px;margin-left: 26%;">Export Dot</el-button>
      </el-main>
    </el-container>
    <el-container id="exportConf">
      <el-header id="eConfH" style="height: 30px;line-height: 30px">Config</el-header>
      <el-main id="eConfM">
        <div >
          <span>fileName : </span>
          <el-input type="text"  v-model="confFileName" style="width: 60%;margin-left: 2%;"></el-input>
        </div>
        <el-row style="margin-top: 20px">
          <template>
            <el-checkbox :indeterminate="isIndeterminate" v-model="checkAllC" @change="handleCheckAllChangeC">全选</el-checkbox>
            <div style="margin: 15px 0;"></div>
            <el-checkbox-group v-model="checkedCitiesC" @change="handleCheckedCitiesChangeC">
              <el-checkbox v-for="city in cities" :label="city" :key="city">{{city}}</el-checkbox>
            </el-checkbox-group>
          </template>
        </el-row>
        <el-button id="findC" @click="ExportConf" style="margin-top: 20px;margin-left: 26%;">ExportConfig</el-button>
      </el-main>
    </el-container>
  </div>


</template>

<script>
  const cityOptions = ['JSON', 'TOML', 'YAML'];
  export default {
    data() {
      return {
        dotFileName:"",
        confFileName:"",
        checkAllC: false,
        checkedCitiesC: ['JSON'],
        cities: cityOptions,
        isIndeterminate: true
      }
    },
    methods: {
      handleCheckAllChangeC(val) {
        this.checkedCitiesC = val ? cityOptions : [];
        this.isIndeterminate = false;
      },
      handleCheckedCitiesChangeC(value) {
        let checkedCount = value.length;
        this.checkAllC = checkedCount === this.cities.length;
        this.isIndeterminate = checkedCount > 0 && checkedCount < this.cities.length;
      },
      ExportDot() {
        alert(this.dotFileName + '.'+this.checkedCities )
      },
      ExportConf() {
        alert(this.confFileName + '.' + this.checkedCitiesC)
      }
    }
  }
</script>
<style>
  #eDotH, #eConfH {
    text-align: left;
    background-color: #D3DCE6;
  }
</style>

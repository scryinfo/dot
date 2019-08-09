<template>
  <div style="line-height: 50px">
    <el-container id="importDot" style="margin-top: -12px;">
      <el-header id="dotH" style="line-height: 30px;height: 30px;text-align: left;">Dot</el-header>
      <el-main id="dotM">
        <div v-for="(item,index) in dotPath">
          <span>dotPath :</span>
          <el-input type="text"  v-model="dotPath[index]" style="width: 60%;margin-left: 2%;"></el-input>
          <el-button id="removeD" @click="del(index)" style="margin-left: 2%">remove</el-button>
        </div>
        <el-row>
          <el-button id="findD" @click="find" style="margin-right:30%;margin-left: 10%;">Import Dot</el-button>
          <!--          <el-button id="addD" @click="add()">Add</el-button>-->
          <!--          <el-button id="removeAllD" @click="removeAll()">Remove All</el-button>-->
        </el-row>
      </el-main>
    </el-container>
    <el-container id="importConf">
      <el-header id="confH" style="height: 30px;line-height: 30px; text-align: left;background-color: #D3DCE6;">Config</el-header>
      <el-main id="confM">
        <div >
          <span>confPath:</span>
          <el-input type="text"  v-model="confPath" style="width: 60%;margin-left: 2%;"></el-input>
          <!--          <el-button id="removeC" @click="del(index)" style="margin-left: 2%">remove</el-button>-->

        </div>
        <el-row>
          <el-button id="findC" @click="ImportConf" style="margin-right:29%;margin-left: 10%;">ImportConfig</el-button>
          <!--          <el-button id="addC" @click="add()">Add</el-button>-->
          <!--          <el-button id="removeAllC" @click="removeAll()">Remove All</el-button>-->
        </el-row>
      </el-main>
    </el-container>
  </div>


</template>

<script>

  export default {
    data() {
      return {
        dotPath:[''],
        confPath:[''],

      }
    },
    methods: {

    }
  }
</script>
<style>

</style>
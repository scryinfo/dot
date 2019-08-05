<template>
  <div id="findDot">
    <el-row>
      <el-button id="altera" @click="ala()">altera</el-button>
      <el-button id="alter" @click="al()">alter</el-button>
      <el-button id="add" @click="add()">Add</el-button>
      <el-button id="removeAll" @click="removeAll()">Remove All</el-button>
    </el-row>
    <div v-for="(item,index) in files">
      <span>filePath{{index+1}}:</span>
      <el-input type="text"  v-model="files[index]" style="width: 78%;margin-left: 2%;"></el-input>
      <el-button id="remove" @click="del(index)" style="margin-left: 2%">remove</el-button>
    </div>
    <el-button id="find" @click="find" style="margin-right:78%;">FindDot</el-button>
  </div>
</template>

<script>
  export default{
    data() {
      return {
        files: [''],
      }
  },
    methods:{
      add(){
        this.files.push('')
      },
      removeAll(){
        this.files=['']
      },
      al(){
        alert(this.files)
      },
      del(index){
        this.files.splice(index,1);
        if (this.files.length==0){
          this.files=['']
        }
      },
      find(){
        
      },
      ala(){
        this.$root.FindDotData.push({name:'bbb',age:19});
        alert(this.$root.FindDotData)
      }
    }
  }
</script>
<style>
  #findDot {
    text-align: right;
    line-height: 50px;
  }
</style>
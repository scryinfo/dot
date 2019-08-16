<template>
  <el-collapse>
    <el-button id="removeAllDot" @click="removeAllDots()" style="margin-bottom: 5px;">Remove All</el-button>
    <el-button id="updateLiveNum" @click="updateLvNum()" style="margin-bottom: 5px;">Update LiveNum</el-button>
    <el-row v-for="(v,index) in this.$root.Dots">
      <el-col :span="3"><div class="grid-content bg-purple" v-text="livesNum[index]" style="text-align: center;line-height: 46px;"></div></el-col>
      <el-col :span="3"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{v.metaData.name}}</div></el-col>
      <el-col :span="10">
        <el-collapse-item v-bind:title="v.metaData.typeId" v-bind:name="index" >
          <el-row v-for="(a,b,c) in v.metaData">
            <el-col :span="6"><div >{{b}}</div></el-col>
            <el-col :span="18"><div >{{a}}</div></el-col>
          </el-row>
        </el-collapse-item>
      </el-col>
      <el-col :span="2"><div style="margin-left: 6px"><el-button v-on:click="open(index)">Name</el-button></div></el-col>
      <el-col :span="3"><div ><el-button v-on:click="addConf(index)">Add Config</el-button></div></el-col>
      <el-col :span="3"> <el-button id="delDot" @click="delDot(index)" style="margin-left: 2%">remove</el-button></el-col>
    </el-row>
  </el-collapse>
</template>
<script>
  export default {
    data() {
      return {
        livesNum: [],
        table: this.$root.Dots,
        activeNames: ['1']
      };
    },
    methods: {
      updateLvNum(){
        this.livesNum=[];
        let l = this.$root.Dots.length;
        for (var i=0;i<l;i++){
          let tId = this.$root.Dots[i].metaData.typeId;
          for (var j=0;j<this.$root.Configs.length;j++){
            if (this.$root.Configs[j].metaData.typeId==tId){
              this.livesNum[i]=this.$root.Configs[j].lives.length;
              break;
            }
          }
        }
      },
      open(index) {
        this.$prompt('输入组件名', 'name', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
        }).then(({value}) => {
          this.$message({
            type: 'success',
            message: '该组件名更改为: ' + value,
          });
          this.$root.Dots[index].metaData.name = value;
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '取消输入'
          });
        });
      },
      addConf(index) {
        let tId = this.$root.Dots[index].metaData.typeId;
        let confNum = this.$root.Configs.length;
        for (var i=0;i<confNum;i++){
          if (this.$root.Configs[i].metaData.typeId == tId){
           alert(" Existed In Config!");
            return
          }
        }
        this.$root.Configs.push(this.$root.Dots[index]);
        alert("Add Config Success!")
      },
      removeAllDots(){
        this.$root.Dots=[]
      },
      delDot(index) {
        this.$root.Dots.splice(index,1)
      },
    }
  }
</script>

<style>
  .bg-purple {
    background: #d3dce6;
  }
  .grid-content {
    border-radius: 4px;
    min-height: 46px;
  }
</style>
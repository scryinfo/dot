<template>
  <div id="findDot" >
    <el-row>
      <el-button id="add" @click="add()">Add</el-button>
      <el-button id="removeAll" @click="removeAll()">Remove All</el-button>
    </el-row>
    <div v-for="(item,index) in files">
      <span>filePath{{index+1}}:</span>
      <el-input type="text"  v-model="files[index]" style="width: 78%;margin-left: 2%;"></el-input>
      <el-button id="remove" @click="del(index)" style="margin-left: 2%">remove</el-button>
    </div>
    <el-button id="find" v-loading="fullscreenLoading" @click="find()" style="margin-right:78%;">FindDot</el-button>
  </div>
</template>

<script>
import {checkType} from "../components/changeDataStructure/checkType";

export default {
    data() {
      return {
        files: [''],
        fullscreenLoading: false,
      }
    },
    methods: {
      add() {
        var end = this.files.length-1;
        if(this.files[end]==''){
          this.$message({
            type:'error',
            message:'the end filepath is empty!'
          })
        }else {
          this.files.push('')
        }
      },
      removeAll() {
        this.files = ['']
      },

      del(index) {
        this.files.splice(index, 1);

      },
      find() {
        var dir = this.files;
        if( dir[0]=='') {
          this.$message({
            type: 'error',
            message: 'Please Input FilePath!'
          })
        } else {
          var {rpcFindDot} = require('../plugins/rpcInterface');
          this.fullscreenLoading = true;
          rpcFindDot(dir, (response) => {
            if (response.getError() != '') {
              var err = response.getError();
              console.log(err);
              this.$message({
                type: 'error',
                message: err,
              });
            } else {
              var res = JSON.parse(response.getDotsinfo());
              var notExitsFile = response.getNoexistdirsList();
              for (var i = 0; i < res.length; i++) {
                var bo = true;
                for (var j = 0, len = this.$root.Dots.length; j < len; j++) {
                  if (res[i].metaData.typeId == this.$root.Dots[j].metaData.typeId) {
                    bo = false;
                    break
                  }
                }
                if (bo) {
                  this.$root.Dots.push(res[i]);
                  this.$root.DotsTem.push(JSON.parse(JSON.stringify(res[i])));
                  this.$root.ExportDots.push(JSON.parse(JSON.stringify(res[i])))
                }
              }
              if (notExitsFile!=[]){
                this.$message({
                  type:'waring',
                  message:'find finish, this not existence Dots:'+notExitsFile
                })
              }else {
                checkType(this.$root.Dots,this.$root.Configs);
                this.$message({
                type: 'success',
                message: 'Find Dot success!'
              })
              }
              this.fullscreenLoading = false;
            }
          });

        }
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

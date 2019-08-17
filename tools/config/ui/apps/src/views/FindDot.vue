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
    <el-button id="find" v-loading="fullscreenLoading" @click="find" style="margin-right:78%;">FindDot</el-button>
  </div>
</template>

<script>

  export default {
    data() {
      return {
        files: [''],
        fullscreenLoading: false,
      }
    },
    methods: {
      add() {
        this.files.push('')
      },
      removeAll() {
        this.files = ['']
      },

      del(index) {
        this.files.splice(index, 1);
        if (this.files.length == 0) {
          this.files = ['']
        }
      },
      find() {
        var dir = this.files;
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
              }
            }
            this.$root.DotsTem = JSON.parse(JSON.stringify(this.$root.Dots));
            this.fullscreenLoading = false;
            this.$message({
              type:'success',
              message:'Find Dot success!'
            })
          }
        });
      },
    }
  }
</script>
<style>
  #findDot {
    text-align: right;
    line-height: 50px;
  }
</style>

<template>
  <el-collapse>
    <el-row v-for="(v,index) in table">
      <el-col :span="3"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">lives num: {{index+1}}</div></el-col>
      <el-col :span="3"><div class="grid-content bg-purple" style="text-align: center;line-height: 46px;">{{v["Meta"]["name"]}}</div></el-col>
      <el-col :span="10">
        <el-collapse-item v-bind:title="v.Meta.typeId" v-bind:name="index" >
          <el-row v-for="(a,b,c) in v.Meta">
            <el-col :span="6"><div >{{b}}</div></el-col>
<!--            <el-col ;span="2"><div>{{c}}</div></el-col>-->
            <el-col :span="18"><div >{{a}}</div></el-col>
          </el-row>
        </el-collapse-item>
      </el-col>
      <el-col :span="2"><div style="margin-left: 6px"><el-button v-on:click="open(index)">Name</el-button></div></el-col>
      <el-col :span="3"><div ><el-button v-on:click="addConf(index)">Add Config</el-button></div></el-col>
    </el-row>
  </el-collapse>
</template>
<script>
  export default {
    data() {
      return {
        table: this.$root.Dots,
        activeNames: ['1']
      };
    },
    methods: {
      open(index) {
        this.$prompt('输入组件名', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
        }).then(({value}) => {
          this.$message({
            type: 'success',
            message: '该组件名更改为: ' + value,
          });
          this.table[index].Meta.name = value;
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '取消输入'
          });
        });
      },
      addConf(index) {
        const h = this.$createElement;
        let n = this.table[index].Meta.name;
        this.$notify({
          title: 'success',
          message: h('i', {style: 'color: teal'},  'add '+ n +' to config success!')
        });
        this.$root.Configs.push(this.table[index]);
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
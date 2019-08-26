<template>
  <el-container>
  <el-header>Dot Config</el-header>
  <el-container>
  <el-aside width="200px">
     <el-menu >
       <router-link to="/findDot" style="text-decoration: none;">
          <el-menu-item index="1">
            <template slot="title"><i class="el-icon-search"></i>Find Dot</template>
          </el-menu-item>
       </router-link>
       <router-link to="/dotList" style="text-decoration: none;">
          <el-menu-item index="2">
            <template slot="title"><i class="el-icon-tickets"></i>Dot List</template>
          </el-menu-item>
       </router-link>
       <router-link to="/config" style="text-decoration: none;">
          <el-menu-item index="3">
            <template slot="title"><i class="el-icon-edit"></i>Config</template>
          </el-menu-item>
       </router-link>
       <router-link to="/import" style="text-decoration: none;">
          <el-menu-item index="4">
            <template slot="title"><i class="el-icon-upload2"></i>Import</template>
          </el-menu-item>
       </router-link>
       <router-link to="/export" style="text-decoration: none;">
          <el-menu-item index="5">
            <template slot="title"><i class="el-icon-download"></i>Export</template>
          </el-menu-item>
       </router-link>
        </el-menu>
  </el-aside>
  <el-main style="height: 90vh;">
    <keep-alive>
      <router-view></router-view>
    </keep-alive>
  </el-main>
  </el-container>
  </el-container>
</template>

<script>

</script>

<style scoped>
    .el-header {
      background-color: #B3C0D1;
      color: #333;
      text-align: center;
      line-height: 60px;
    }
  .el-aside {
    background-color: #D3DCE6;
    color: #333;
    text-align: center;
    line-height: 200px;
  }

  .el-main {
    background-color: #E9EEF3;
    color: #333;
  }
.el-menu-item {
  background: #D3DCE6;
}
  body > .el-container {
    margin-bottom: 40px;
  }

</style>

const plugins = [
    "@vue/babel-plugin-transform-vue-jsx",
  [
    '@babel/plugin-transform-modules-commonjs',
    {
      allowTopLevelThis: true
    }
  ]
]
if(process.env.NODE_ENV === 'production') {
  plugins.push("transform-remove-console")
}
module.exports = {
  presets: [
    '@vue/app'
  ],
  plugins: plugins
}

const path = require('path')
module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  chainWebpack: config => {
    config
      .entry('app')
      .clear()
      .add('./src/main.ts')
      .end();
    config.module.rule('ts')
      .test(/\.ts$/)
      .use('ts-loader')
      .loader('ts-loader')
      .options({appendTsSuffixTo: [/\.vue$/]});
  },
  
  devServer: {

    host: '0.0.0.0'
    }
}

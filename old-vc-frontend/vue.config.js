module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  chainWebpack: config => {
    config.resolve.alias.set('vue', '@vue/compat')
    config.module
      .rule('vue')
      .use('vue-loader')
      .tap((options) => {
        return {
          ...options,
          compilerOptions: {
            compatConfig: {
              MODE: 2
            }
          }
        }
      })
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

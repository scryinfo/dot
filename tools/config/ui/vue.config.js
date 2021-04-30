const fs = require("fs");
let path = require("path");
let glob = require("glob");
const CopyPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CompressionWebpackPlugin = require("compression-webpack-plugin");
const productionGzipExtensions = /\.(js|css|txt|html|ico|svg)(\?.*)?$/i; // del 'json' type
const merge = require("webpack-merge");

const pages = getEntry("src/views/**/main.ts");

function getEntry(globPath) {
    let pages = {};
    glob.sync(globPath).forEach(function (entry) {
        let dname = path.dirname(entry);
        let basename = path.basename(dname);
        let template = path.join(dname, "index.html");
        if (!fs.existsSync(template)) {
            template = "scr/views/index.html";
        }

        pages[basename] = {
            entry: entry,
            template: template,
            filename: `${basename}.html`
            // chunks: [basename]
        };

        // console.log(pages[basename]);
    });
    return pages;
}

module.exports = {
    publicPath: "./",
    assetsDir: "assets", // 打包后静态资源的位置
    productionSourceMap: false,
    configureWebpack: config => {
        config.output.filename = "assets/js/[name].[hash:8].js";
        config.output.chunkFilename = "assets/js/[name].[hash:8].js";
        config.plugins.push(new CopyPlugin([{from: "src/assets", to: "assets"}]));

        if (process.env.NODE_ENV === "production") {
            const plugins = [];
            plugins.push(
                new CompressionWebpackPlugin({
                    filename: "[path].gz[query]",
                    algorithm: "gzip",
                    test: productionGzipExtensions,
                    threshold: 10240,
                    minRatio: 0.8
                })
            );
            plugins.push(
                new CompressionWebpackPlugin({
                    filename: "[path].br[query]",
                    algorithm: "brotliCompress",
                    test: /\.(js|css|html|svg)$/,
                    compressionOptions: {level: 11},
                    threshold: 10240,
                    minRatio: 0.8
                })
            );
            config.plugins = [...config.plugins, ...plugins];
            config.optimization = merge(config.optimization, {
                minimize: true
            });
        } else {
            // mutate for development...
        }
    },
    chainWebpack: config => {
        config.resolve.alias.set("@", path.join(__dirname, "src"));
        config.resolve.alias.set("vue$", "vue/dist/vue.runtime.esm.js");

        config.plugin("extract-css").use(
            new MiniCssExtractPlugin({
                filename: "assets/css/[name].[hash:8].css",
                chunkFilename: "assets/css/[name].[hash:8].css"
            })
        );

        config.module
            .rule("images")
            .test(/\.(gif|png|jpe?g|svg)$/i)
            .use("url-loader")
            .loader("file-loader")
            .options({
                name: "assets/img/[name].[ext]"
            });

        if (process.env.NODE_ENV === "production") {
            //
        } else {
            //
        }
    },
    pages: pages,
    css: {
        // 是否使用css分离插件 ExtractTextPlugin
        extract: false,
        // 开启 CSS source maps?
        sourceMap: false,
        // css预设器配置项
        loaderOptions: {},
        // 启用 CSS modules for all css / pre-processor files.
        requireModuleExtension: true
    },
    devServer: {
        contentBase: "./dist",
        host: "",
        port: "9689",
        open: true, //打开浏览器
        openPage: "/home.html",
        hot: true, //热加载,
        clientLogLevel: "debug",
        overlay: {
            warnings: false,
            errors: false
        }
    },
    lintOnSave: process.env.NODE_ENV === "production"
};

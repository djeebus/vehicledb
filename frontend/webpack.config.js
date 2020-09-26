const path = require('path')
const MiniCssExtractPlugin = require("mini-css-extract-plugin")
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const TerserWebpackPlugin = require('terser-webpack-plugin')
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin')

module.exports = function (_env, argv) {
    const isProduction = argv.mode === "production"
    const isDevelopment = !isProduction

    return {
        devtool: isDevelopment && "cheap-module-source-map",
        devServer: {
            historyApiFallback: true,
        },
        entry: "./src/index.js",
        output: {
            path: path.resolve(__dirname, "dist"),
            filename: "assets/js/[name].[contenthash:8].js",
            publicPath: "/"
        },
        module: {
            rules: [
                {
                    test: /\.jsx?$/,
                    exclude: /node_modules/,
                    use: {
                        loader: "babel-loader",
                        options: {
                            cacheDirectory: true,
                            cacheCompression: false,
                            envName: isProduction ? "production" : "development",
                        },
                    },
                },
                {
                    test: /\.css$/,
                    use: [
                        isProduction ? MiniCssExtractPlugin.loader : "style-loader",
                        "css-loader",
                    ],
                },
            ],
        },
        resolve: {
            extensions: [".js", ".jsx"]
        },
        plugins: [
            isProduction && new MiniCssExtractPlugin({
                filename: "assets/css/[name].[contenthash:8].css",
                chunkFilename: "assets/css/[name].[contenthash:8].chunk.css",
            }),
            new webpack.DefinePlugin({
                "process.env.NODE_ENV": JSON.stringify(
                    isProduction ? "production" : "development",
                ),
                "process.env.API_URL": JSON.stringify(
                    isProduction ? "https://api.vehicledb.com" : "http://localhost:8000",
                ),
            }),
            new HtmlWebpackPlugin({
                template: path.resolve(__dirname, "public/index.html"),
                inject: true,
            }),
        ].filter(Boolean),
        optimization: {
            minimize: isProduction,
            minimizer: [
                new TerserWebpackPlugin({
                    terserOptions: {
                        compress: {comparisons: false},
                        mangle: {safari10: true},
                        output: {
                            comments: false,
                            ascii_only: true,
                        },
                        warnings: false,
                    },
                }),
                new OptimizeCssAssetsPlugin(),
            ],
        },
    }
}
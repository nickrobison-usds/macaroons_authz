const path = require('path');
const CleanWebpackPlugin = require('clean-webpack-plugin');
const webpack = require('webpack');

const outputPath = path.join(__dirname, '..', 'dist');


module.exports = {
    entry: ["./src/app.ts"],
    module: {
        rules: [
            {
                test: /.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            }
        ]
    },
    target: 'node',
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
    },
    plugins: [
        new CleanWebpackPlugin(['dist'], {
            root: path.join(__dirname, '..'),
            verbose: true
        }),
        // We need this in order to bundle postgres natively.
        new webpack.IgnorePlugin(/^pg-native$/)
    ],
    output: {
        path: outputPath,
        filename: 'target_service.js'
    }
};

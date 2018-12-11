const path = require('path');
const CleanWebpackPlugin = require('clean-webpack-plugin');

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
        })
    ],
    output: {
        path: outputPath,
        filename: 'target_service.js'
    }
};

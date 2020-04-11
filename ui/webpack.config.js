const VueLoaderPlugin = require('vue-loader/lib/plugin')
const exec = require("child_process").exec;

module.exports = {
    entry: "./src/sermoni.js",
    output: {
        filename: "sermoni.js"
    },
    module: {
        rules: [
            {test: /\.js$/, use: 'babel-loader'},
            {test: /\.vue$/, use: 'vue-loader'},
            {test: /\.css$/, use: ['vue-style-loader', 'css-loader']},
        ]
    },
    plugins: [
        new VueLoaderPlugin(),
        // Ad-hoc plugin for running script automatically
        // Thanks to https://stackoverflow.com/a/49786887
        {
            apply: (compiler) => {
                compiler.hooks.afterEmit.tap("AfterEmitPlugin", (compilation) => {
                    exec("./generate.sh", (err, stdout, stderr) => {
                        if (stdout) process.stdout.write(stdout);
                        if (stderr) process.stderr.write(stderr);
                    });
                });
            }
        }
    ]
}

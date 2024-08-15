import { nodeResolve } from "@rollup/plugin-node-resolve";
import commonjs from "@rollup/plugin-commonjs";
import json from "@rollup/plugin-json";

export default {
    input: "dist/index.js",
    output: {
        file: "dist/bundle.js",
        format: "cjs",
    },
    plugins: [
        nodeResolve({
            preferBuiltins: true,
        }),
        commonjs(),
        json(),
    ],
};

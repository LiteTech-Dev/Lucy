import path from "path";
import { FileUtil } from "../util/fileUtil.js";

try {
    var res = await FileUtil.ExtractFromJar("EFMCore-1.0.jar", "plugin.yml");
    console.log(res);
} catch (err) {
    console.log(err)
}
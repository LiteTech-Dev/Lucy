import * as jszip from "jszip";
import * as fs from "fs";
import { getMCPMDir } from "./directory.js";
import path from "path";

export async function extractFromJar(
    jarPath: string,
    target: string
): Promise<fs.ReadStream | null> {
    const extractDir: path.ParsedPath = path.parse(
        path.join(getMCPMDir().dir, path.resolve("extracted"))
    );

    let jarStructure: string[] | null = await getJarStructure(jarPath);
    if (jarStructure === null) return null;

    jarStructure.forEach((element) => {
        if (element === target) {
            // 在这里插入解压并写入.mcpm的部分
            // 并且返回对该文件的fs.readStream
        }
    });

    return null;
}

export async function getJarStructure(
    jarPath: string
): Promise<string[] | null> {
    try {
        const jar = await jszip.loadAsync(jarPath);
        return Object.keys(jar.files);
    } catch (err) {
        console.error(" ", err); // 本地化
    }
    return null;
}

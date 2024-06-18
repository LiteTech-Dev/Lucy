import * as jszip from "jszip";
import * as fs from "fs";

export async function extractFromJar(
    jarPath: string,
    target: string
): Promise<fs.ReadStream | null> {
    let jarStructure: string[] | null = await getJarStructure(jarPath);
    if (jarStructure === null) return null;

    jarStructure.forEach((element) => {
        if (element === target) {
            // 在这里插入解压并写入.mcpm的部分
            // return;
            // 返回对该文件的fs.readStream
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

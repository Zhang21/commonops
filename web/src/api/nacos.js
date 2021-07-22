import req from "../utils/axios";

const postNacosServer = (params) => {
    return req.post("configCenter/nacos", params);
};
export { postNacosServer };

const getNacosList = (params) => {
    return req.get("configCenter/nacos/list", params);
};
export { getNacosList };

const getNacosNamespaceList = (params) => {
    return req.get("configCenter/nacos/namespaces", params);
};
export { getNacosNamespaceList };

const getNacosConfigsList = (params) => {
    return req.get("configCenter/nacos/configs", params);
};
export { getNacosConfigsList };

const postNacosConfig = (params) => {
    return req.post("configCenter/nacos/config", params);
};
export { postNacosConfig };

const putNacosConfig = (params) => {
    return req.put("configCenter/nacos/config", params);
};
export { putNacosConfig };

const deleteNacosConfig = (params) => {
    return req.delete("configCenter/nacos/config", params);
};
export { deleteNacosConfig };

const postNacosConfigCopy = (params) => {
    return req.post("configCenter/nacos/config/copy", params);
};
export { postNacosConfigCopy };
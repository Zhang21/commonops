import React, { Component, Fragment } from "react";
import {
    Layout,
    Row,
    message,
    Button,
    Col,
    Typography,
    Select,
    Form,
    Drawer,
    Input,
    Radio, Table, Modal, Divider,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import { withRouter } from 'react-router-dom';
import { getNacosList, getNacosNamespaceList, getNacosConfigsList,
    putNacosConfig, postNacosConfig, postNacosServer, deleteNacosConfig, postNacosConfigCopy } from "../../api/nacos";

const { Paragraph, Text } = Typography;
const { Content } = Layout;
const { Option } = Select;
const { TextArea } = Input;

class NacosServerForm extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        };
        return (
            <Form ref={this.props.formRef} {...formItemLayout}>
                <Form.Item
                    label="集群名称"
                    name="alias"
                    rules={[{ required: true, message: "集群名称不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="访问地址"
                    name="endpoint"
                    rules={[{ required: true, message: "endpoint不能为空！" }]}
                >
                    <Input addonBefore="http://" />
                </Form.Item>
                <Form.Item
                    label="用户名"
                    name="username"
                    rules={[{ required: true, message: "username不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="密码"
                    name="password"
                    rules={[{ required: true, message: "password不能为空！" }]}
                >
                    <Input />
                </Form.Item>
            </Form>
        );
    }
}

class NacosConfigForm extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 16 },
        };
        return (
            <Form ref={this.props.formRef} {...formItemLayout}>
                <Form.Item
                    label="ConfigId"
                    name="configId"
                    hidden={true}
                >
                </Form.Item>
                <Form.Item
                    label="DataId"
                    name="dataId"
                    rules={[{ required: true, message: "DataId不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Group"
                    name="group"
                    rules={[{ required: true, message: "Group不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="格式"
                    name="configType"
                    rules={[{ required: true, message: "格式不能为空！" }]}
                >
                    <Radio.Group>
                        <Radio value="yaml">YAML</Radio>
                        <Radio value="properties">Properties</Radio>
                        <Radio value="text">TEXT</Radio>
                        <Radio value="json">JSON</Radio>
                        <Radio value="xml">XML</Radio>
                        <Radio value="html">HTML</Radio>
                    </Radio.Group>
                </Form.Item>
                <Form.Item
                    label="配置内容"
                    name="content"
                    rules={[{ required: true, message: "配置内容不能为空！" }]}
                >
                    <TextArea rows={6} allowClear/>
                </Form.Item>
            </Form>
        );
    }
}

class NacosConfigCopyForm extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 15 },
        };
        return (
            <Form ref={this.props.formRef} {...formItemLayout}>
                <Form.Item
                    label="srcNamespace"
                    name="srcNamespace"
                    hidden={true}
                />
                <Form.Item
                    label="srcDataId"
                    name="srcDataId"
                    hidden={true}
                />
                <Form.Item
                    label="srcGroup"
                    name="srcGroup"
                    hidden={true}
                />
                <Form.Item
                    label="Namespace"
                    name="dstNamespace"
                    rules={[{ required: true, message: "Namespace不能为空！" }]}
                >
                    <Select>
                        {this.props.nsList.map((item, index)=>{
                            return <Option key={index} value={item.NamespaceShowName}>{item.NamespaceShowName}</Option>
                        })}
                    </Select>
                </Form.Item>
                <Form.Item
                    label="DataId"
                    name="dstDataId"
                    rules={[{ required: true, message: "DataId不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Group"
                    name="dstGroup"
                    rules={[{ required: true, message: "Group不能为空！" }]}
                >
                    <Input />
                </Form.Item>
            </Form>
        );
    }
}

class NacosContent extends Component {
    constructor(props) {
        super(props);
        this.configFormRef = React.createRef();
        this.configCopyFormRef = React.createRef();
        this.serverFormRef = React.createRef();
        this.state = {
            addNacosModalVisible: false,
            addNacosConfigModalVisible: false,
            copyNacosConfigModalVisible: false,
            currentNacosServerId: "",
            currentNs: "",
            nacosList: [],
            nsList: [],
            columns: [
                {
                    title: "DataId",
                    dataIndex: "dataId",
                    key: "dataId",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "Group",
                    dataIndex: "group",
                    key: "group",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "配置类型",
                    dataIndex: "type",
                    key: "type",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "操作",
                    dataIndex: "操作",
                    key: "操作",
                    width: 300,
                    render: (value, record) => {
                        return (
                            <div style={{textAlign: "center"}}>
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showConfigContent.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    详情
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.updateConfigContent.bind(
                                        this,
                                        record,
                                    )}
                                    disabled={!this.props.aclAuthMap["PUT:/configCenter/nacos/config"]}
                                >
                                    修改
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.copyConfigContent.bind(
                                        this,
                                        record,
                                    )}
                                    disabled={!this.props.aclAuthMap["POST:/configCenter/nacos/config/copy"]}
                                >
                                    复制
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="danger"
                                    size="small"
                                    onClick={this.deleteConfigContent.bind(
                                        this,
                                        record,
                                    )}
                                    disabled={!this.props.aclAuthMap["DELETE:/configCenter/nacos/config"]}
                                >
                                    删除
                                </Button>
                            </div>
                        )
                    },
                },
            ],
            tableData: [],
            pagination: {
                showSizeChanger: true,
                pageSizeOptions: ["10", "20"],
                onShowSizeChange: (current, size) =>
                    this.onShowSizeChange(current, size),
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (page, pageSize) => this.changePage(page, pageSize),
            },
            tableLoading: false,
            configContentDrawerVisible: false,
        };
    }

    componentDidMount() {
        this.loadNacosServerList();
    }

    loadNacosServerList() {
        getNacosList().then((res)=>{
            if(res.code===0) {
                this.setState({nacosList: res.data});
            } else {
                message.error(res.msg);
            }
        })
    }

    loadNacosNamespaceList() {
        if(this.state.currentNacosServerId!=="") {
            getNacosNamespaceList({id: this.state.currentNacosServerId}).then((res)=>{
                if(res.code===0) {
                    this.setState({nsList: res.data.data});
                } else {
                    message.error(res.msg);
                }
            })
        }
    }

    loadNamespaceConfigs() {
        if(this.state.currentNacosServerId==="") {
            message.info("请选择集群!");
            return
        }
        if(this.state.currentNs==="") {
            message.info("请选择命名空间!");
            return
        }
        this.setState({ tableLoading: true });
        const queryParams = {
            page: parseInt(this.state.pagination.page),
            size: parseInt(this.state.pagination.pageSize),
            id: this.state.currentNacosServerId,
            namespace: this.state.currentNs
        };
        getNacosConfigsList(queryParams).then((res)=>{
            if(res.code===0) {
                var data = JSON.parse(res.data);
                const pagination = this.state.pagination;
                pagination.total = parseInt(data.totalCount);
                pagination.page = parseInt(data.pageNumber);
                pagination.showTotal(parseInt(data.totalCount));
                this.setState({
                    pagination,
                });
                this.setState({tableData: data.pageItems})
            } else {
                message.error(res.msg);
            }
            this.setState({ tableLoading: false });
        })
    }

    onShowSizeChange(current, size) {
        let pagination = {
            ...this.state.pagination,
            page: 1,
            current: 1,
            pageSize: size,
        };
        this.setState(
            {
                pagination: pagination,
            },
            () => {
                this.loadNamespaceConfigs();
            },
        );
    }

    changePage = (page, pageSize) => {
        this.setState(
            {
                pagination: {
                    ...this.state.pagination,
                    page: page,
                    current: page,
                    pageSize: pageSize,
                },
            },
            () => {
                this.loadNamespaceConfigs();
            },
        );
    };


    showConfigContent(record) {
        this.setState({
            configContentDrawerVisible: true,
            configType: record.type,
            configContent: record.content,
        });
    }

    onCloseConfigContentDrawer = () => {
        this.setState({
            configContentDrawerVisible: false,
        });
    }

    copyConfigContent(record) {
        let that = this;
        this.setState({copyNacosConfigModalVisible: true}, ()=>{
            setTimeout(()=>{
                that.configCopyFormRef.current.setFieldsValue({
                    srcNamespace: record.tenant,
                    srcDataId: record.dataId,
                    srcGroup: record.group,
                });
            }, 300);
        })

    }

    updateConfigContent(record) {
        let that = this;
        this.setState({addNacosConfigModalVisible: true}, ()=>{
            setTimeout(()=>{
                that.configFormRef.current.setFieldsValue({
                    configId: record.id,
                    dataId: record.dataId,
                    group: record.group,
                    configType: record.type,
                    content: record.content,
                });
            }, 200);
        });
    }

    deleteConfigContent(record) {
        deleteNacosConfig({
            id: ""+this.state.currentNacosServerId,
            namespace: this.state.currentNs,
            dataId: record.dataId,
            group: record.group,
        }).then((res)=>{
            if(res.code===0) {
                message.success("删除成功");
                this.loadNamespaceConfigs();
            } else {
                message.error(res.msg);
            }
        })
    }

    addNacosServer() {
        this.setState({addNacosModalVisible: true});
    }

    addNacosConfig() {
        if(this.state.currentNacosServerId==="") {
            message.info("请选择集群!");
            return
        }
        if(this.state.currentNs==="") {
            message.info("请选择命名空间!");
            return
        }
        this.setState({addNacosConfigModalVisible: true});
    }

    handleNacosChange(e) {
        this.setState({currentNacosServerId: e, currentNs: ""}, ()=>{
            this.loadNacosNamespaceList();
        });
    }

    handleNsChange(e) {
        this.setState({currentNs: e}, ()=>{
            this.loadNamespaceConfigs();
        });
    }

    nacosConfigQuery() {
        this.loadNamespaceConfigs();
    }

    submitAddNacosServer = (e) => {
        e.preventDefault();
        this.serverFormRef.current.validateFields().then((values) => {
            let params = {
                ...values,
            };
            postNacosServer(params).then((res)=>{
                if(res.code===0) {
                    this.setState({addNacosModalVisible: false});
                    this.loadNacosServerList();
                } else {
                    message.error(res.msg);
                }
            })
        });
    };

    cancelAddNacosServer() {
        this.setState({addNacosModalVisible: false});
    }

    submitCreateConfig = (e) => {
        e.preventDefault();
        this.configFormRef.current.validateFields().then((values) => {
            let params = {
                ...values,
                id:  ""+this.state.currentNacosServerId,
                namespace: this.state.currentNs,
            };
            if(values.configId===undefined) {
                postNacosConfig(params).then((res)=>{
                    if(res.code===0) {
                        message.success("创建成功")
                        this.setState({addNacosConfigModalVisible: false});
                        this.loadNamespaceConfigs();
                    } else {
                        message.error(res.msg);
                    }
                })
            } else {
                putNacosConfig(params).then((res) => {
                    if (res.code === 0) {
                        message.success("修改成功")
                        this.setState({addNacosConfigModalVisible: false});
                        this.loadNamespaceConfigs();
                    } else {
                        message.error(res.msg);
                    }
                })
            }
        });
    };

    submitCreateConfigCopy = (e) => {
        e.preventDefault();
        this.configCopyFormRef.current.validateFields().then((values) => {
            let params = {
                ...values,
                id: ""+this.state.currentNacosServerId,
            };
            postNacosConfigCopy(params).then((res)=>{
                if(res.code===0) {
                    message.success("复制成功!");
                    this.setState({copyNacosConfigModalVisible: false});
                } else {
                    message.error(res.msg);
                }
            });
        });
    };

    cancelCreateConfig() {
        this.setState({addNacosConfigModalVisible: false});
    }

    cancelCreateConfigCopy() {
        this.setState({copyNacosConfigModalVisible: false});
    }

    render() {
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: "5px 20px",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["配置中心", "Nacos管理"]} />
                <Row style={{marginBottom: 20}}>
                    <Col span={7}>
                        选择集群:&nbsp;&nbsp;
                        <Select style={{ width: "200px" }} onChange={this.handleNacosChange.bind(this)}>
                            {this.state.nacosList.map((item, index)=>{
                                return <Option key={index} value={item.Id}>{item.Alias}</Option>
                            })}
                        </Select>
                    </Col>
                    <Col span={7}>
                        选择命名空间:&nbsp;&nbsp;
                        <Select style={{ width: "200px" }} onChange={this.handleNsChange.bind(this)}>
                            {this.state.nsList.map((item, index)=>{
                                return <Option key={index} value={item.NamespaceShowName}>{item.NamespaceShowName}</Option>
                            })}
                        </Select>
                    </Col>
                    <Col span={10}>
                        <Button
                            onClick={this.nacosConfigQuery.bind(this)}
                            style={{ width: "100px" }}
                            disabled={!this.props.aclAuthMap["GET:/configCenter/nacos/configs"]}
                        >
                            查询
                        </Button> &nbsp;&nbsp;
                        <Button
                            type="primary"
                            onClick={this.addNacosConfig.bind(this)}
                            style={{ width: "130px" }}
                            disabled={!this.props.aclAuthMap["POST:/configCenter/nacos/config"]}
                        >
                            新增配置
                        </Button> &nbsp;&nbsp;
                        <Button
                            type="primary"
                            onClick={this.addNacosServer.bind(this)}
                            style={{ width: "130px" }}
                            disabled={!this.props.aclAuthMap["POST:/configCenter/nacos"]}
                        >
                            添加Nacos集群
                        </Button>
                    </Col>
                </Row>

                <Modal
                    title="新增Nacos集群"
                    visible={this.state.addNacosModalVisible}
                    onOk={this.submitAddNacosServer}
                    onCancel={this.cancelAddNacosServer.bind(this)}
                    destroyOnClose={true}
                >
                    <NacosServerForm formRef={this.serverFormRef}/>
                </Modal>

                <Modal
                    title="配置信息"
                    visible={this.state.addNacosConfigModalVisible}
                    onOk={this.submitCreateConfig}
                    onCancel={this.cancelCreateConfig.bind(this)}
                    width={800}
                    destroyOnClose={true}
                >
                    <NacosConfigForm formRef={this.configFormRef} />
                </Modal>

                <Modal
                    title="配置复制"
                    visible={this.state.copyNacosConfigModalVisible}
                    onOk={this.submitCreateConfigCopy}
                    onCancel={this.cancelCreateConfigCopy.bind(this)}
                    width={500}
                    destroyOnClose={true}
                >
                    <NacosConfigCopyForm nsList={this.state.nsList} formRef={this.configCopyFormRef} />
                </Modal>

                <Drawer
                    title="配置信息"
                    placement="left"
                    width={800}
                    closable={false}
                    onClose={this.onCloseConfigContentDrawer}
                    visible={this.state.configContentDrawerVisible}
                >
                    <Text strong>配置类型: {this.state.configType}</Text>
                    <Divider />
                    <Text strong>配置内容: </Text>
                    <Paragraph style={{display: "inline-block", width: "40px"}} copyable={{icon: <span>复制</span>, text:this.state.configContent}}/>
                    <pre style={{backgroundColor: "rgb(36, 35, 35)", color: "#eee", padding: "10px 10px"}}>
                        {this.state.configContent}
                    </pre>
                </Drawer>

                <Table
                    columns={this.state.columns}
                    dataSource={this.state.tableData}
                    scroll={{ x: "max-content" }}
                    pagination={this.state.pagination}
                    loading={this.state.tableLoading}
                    bordered
                    size="small"
                />
            </Content>
        );
    }
}

export default withRouter(NacosContent);

<template>
  <el-row :gutter="20" v-loading="loading">
    <el-tabs tab-position="left">
      <el-tab-pane label="代理配置">
        <el-col :span="8">
          <el-table :data="server.info.data" size="mini" :stripe="true" :show-header="false" style="width: 100%">
            <el-table-column prop="name" label="配置" width="60"></el-table-column>
            <el-table-column prop="value" label="值" width="110">
              <template slot-scope="scope">
                <span v-if="scope.row.key === 'status'">
                  <el-tag v-if="scope.row.value" type="success" size="mini" effect="plain">运行中</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">停止</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'balance'">
                  <el-tag v-if="!scope.row.value" type="success" size="mini" effect="plain">随机模式</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">循环模式</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'all_proxy'">
                  <el-tag v-if="scope.row.value" type="success" size="mini" effect="plain">全部转发</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">规则转发</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'instances'">
                  <el-tag v-for="(item, i) in scope.row.value" size="mini" type="danger" effect="plain" :key="i">{{item.address}}</el-tag>
                </span>
                <span v-else>{{scope.row.value}}</span>
              </template>
            </el-table-column>
            <el-table-column prop="operation" align="right">
              <template slot-scope="scope">
                <span v-if="scope.row.key === 'status'">
                  <el-button v-if="scope.row.value" @click="stop_server" type="text" size="small">停止</el-button>
                  <el-button v-else @click="start_server" type="text" size="small">启动</el-button>
                </span>
                <span v-else-if="scope.row.key === 'port'">
                  <el-button @click.native="edit_server(scope.row)" type="text" size="small">编辑</el-button>
                </span>
              </template>
            </el-table-column>
          </el-table>

          <el-card v-if="server.info.instances.length > 0" style="margin-top: 10px;">
            <div slot="header" class="clearfix">
              <span>代理池</span>
            </div>
            <div v-for="(item, i) in server.info.instances" :key="i">
              <el-tag  size="mini" effect="plain" >{{item.address}}</el-tag>
              <el-popconfirm title="是否将实例从代理池中移除？" @onConfirm="instance_remove_from_pool(item.id)">
                <el-button slot="reference" style="color: red;float: right;" type="text" size="mini" icon="el-icon-remove"></el-button>
              </el-popconfirm>
              <el-divider></el-divider>
            </div>
          </el-card>
        </el-col>
        <!-- // proxy instances  -->
        <el-col :span="16">
          <el-table :data="instanceData" :stripe="true" size="mini" style="width: 100%">
            <el-table-column prop="address">
              <template slot="header" slot-scope="scope">
                <el-button type="text" size="mini" icon="el-icon-circle-plus-outline" @click="instaces.create.visit = true">新增实例</el-button>
              </template>
              <template slot-scope="scope">
                  <el-button type="text" size="mini" @click="edit_instance(scope.row)">
                    {{`ssh://${scope.row.username}@${scope.row.address}`}}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="delay" label="延迟" width="100">
              <template slot-scope="scope" >
                <el-tag v-if="scope.row.delay < 100" size="mini" type="success" effect="plain">{{scope.row.delay}} ms</el-tag>
                <el-tag v-else-if="scope.row.delay < 200" size="mini" type="warning" effect="plain">{{scope.row.delay}} ms</el-tag>
                <el-tag v-else size="mini" type="danger" effect="plain">{{scope.row.delay}} ms</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template slot-scope="scope" >
                <el-popconfirm title="是否将实例删除？" @onConfirm="instance_remove(scope.row.id)">
                  <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-error"></el-button>
                </el-popconfirm>
                <el-popconfirm v-if="!scope.row.status" title="是否将实例加入代理池？" @onConfirm="instance_into_pool(scope.row.id)">
                  <el-button slot="reference" style="color: green" type="text" size="mini" icon="el-icon-circle-plus"></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-col>
      </el-tab-pane>
      <el-tab-pane label="规则配置">
        <el-button-group>
          <el-button size="mini" type="primary" icon="el-icon-edit" @click="roles.create.visit = true">新增规则</el-button>
        </el-button-group>
        <!-- <el-input size="mini" placeholder="请输入内容" v-model="roles.filter" style="width: 250px;float: right;">
          <el-button slot="append" icon="el-icon-search"></el-button>
        </el-input> -->

        <el-table :data="rolesData" stripe border size="mini" style="margin-top: 10px;">
          <el-table-column prop="sub" label="域名" width="200">
            <template slot-scope="scope">
              {{`${scope.row.sub}.${scope.row.domain}`}}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="规则"
          :filters="[{ text: '代理', value: true }, { text: '封禁', value: false }]"
          :filter-method="filter_status"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.status">
                <el-tag effect="plain" size="mini">转发 -> {{get_instance(scope.row.instance_id)}}</el-tag>
              </span>
              <el-tag v-else size="mini" type="danger" effect="plain">封禁</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" fixed="right" width="100">
            <template slot-scope="scope">
              <el-popconfirm title="是否删除该规则？" @onConfirm="remove_role(scope.row)">
                <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-delete"></el-button>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog title="修改服务" :visible.sync="server.edit.visit">
      <el-form :model="server.edit.form" label-position="right">
        <el-form-item label="端口" label-width="100px">
          <el-input-number :min="81" :max="65534" controls-position="right" v-model="server.edit.form.port" size="small" autocomplete="off"></el-input-number>
        </el-form-item>
        <el-form-item label="负载" label-width="100px">
          <el-switch
            v-model="server.edit.form.status"
            active-text="轮训"
            inactive-text="随机">
          </el-switch>
        </el-form-item>
        <el-form-item label="模式" label-width="100px">
          <el-switch
            v-model="server.edit.form.all_proxy"
            active-text="全部转发"
            inactive-text="规则转发">
          </el-switch>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="server.edit.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_edit_server">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="新增实例" :visible.sync="instaces.create.visit">
      <el-form :model="instaces.create.form" label-position="right">
        <el-form-item label="地址" label-width="100px">
          <el-input v-model="instaces.create.form.address" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="用户" label-width="100px">
          <el-input v-model="instaces.create.form.username" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" label-width="100px">
          <el-input v-model="instaces.create.form.password" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="秘钥" label-width="100px">
          <el-input type="textarea" :autosize="{ minRows: 3, maxRows: 6}" placeholder="请输入内容" v-model="instaces.create.form.private_key"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="instaces.create.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_create_instance">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="修改实例" :visible.sync="instaces.edit.visit">
      <el-form :model="instaces.edit.form" label-position="right">
        <el-form-item style="display: none" label="ID" label-width="100px">
          <el-input v-model="instaces.edit.form.ID" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="地址" label-width="100px">
          <el-input v-model="instaces.edit.form.address" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="用户" label-width="100px">
          <el-input v-model="instaces.edit.form.username" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" label-width="100px">
          <el-input v-model="instaces.edit.form.password" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="秘钥" label-width="100px">
          <el-input type="textarea" :autosize="{ minRows: 3, maxRows: 6}" placeholder="请输入内容" v-model="instaces.edit.form.private_key"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="instaces.edit.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_edit_instance">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="新增规则" :visible.sync="roles.create.visit">
      <el-form :model="roles.create.form" label-position="right">
        <el-form-item label="域名" label-width="100px">
          <el-input v-model="roles.create.form.url" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="规则" label-width="100px">
          <el-switch v-model="roles.create.form.status" active-text="转发" inactive-text="封禁" @change="role_change"></el-switch>
        </el-form-item>
        <el-form-item label="实例" label-width="100px" :style="roles.create.style">
          <el-select v-model="roles.create.form.instance_id" size="small" placeholder="请选择转发">
            <el-option label="默认转发" value=""></el-option>
            <el-option v-for="(item, i) in instanceData" :key="i" :label="item.address" :value="item.id"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="roles.create.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_create_role">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>

<script>

export default {
  data() {
    return {
      loading: false,
      instanceData: [],
      serverData: [],
      proxy: {
        Status: false,
        Error: null
      },
      server: {
        edit: {
          visit: false,
          form: {
            name: "",
            port: 0,
            status: false,
            username: "",
            password: "",
            auto_proxy: false,
            all_proxy: false
          }
        },
        info: {
          data: [],
          instances: []
        }
      },
      instaces: {
        create: {
          visit: false,
          form: {
            address: undefined,
            username: undefined,
            password: undefined,
            privateKey: undefined
          },
        },
        edit: {
          visit: false,
          form: {
            id: undefined,
            address: undefined,
            username: undefined,
            password: undefined,
            private_key: undefined
          },
        }
      },
      rolesData: [],
      instanceData: [],
      roles: {
        filter: undefined,
        create: {
          visit: false,
          style: "display: none",
          form: {
            url: undefined,
            status: undefined,
            instance_id: ""
          },
        }
      }
    }
  },
  methods: {
    edit_instance (item) {
      this.instaces.edit.visit = true
      this.instaces.edit.form = item
    },
    edit_server (item) {
      let that = this
      this.server.info.data.forEach(function (row) {
        switch (row.key) {
          case "name":
            that.server.edit.form.name = row.value
          case "port":
            that.server.edit.form.port = row.value
          case "balance":
            that.server.edit.form.status = row.value
          case "all_proxy":
            that.server.edit.form.all_proxy = row.value
        }
      })
      this.server.edit.visit = true
    },
    submit_create_instance: function () {
      let that = this
      this.instaces.create.visit = false
      this.$api.post("/proxy/instance/create", this.instaces.create.form).then(function (response) {
        that.$message({type: "success", message: '创建成功'})
        that.refresh_instances()
      }).catch(function (response) {
        that.$message({type: "error", message: response.message})
      })
    },
    submit_edit_server () {
      let that = this
      this.$api.post('/proxy/server/update', this.server.edit.form).then(function (response) {
        that.server.edit.visit = false
        that.$message({message: '修改成功', type: 'success'})
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    submit_edit_instance () {
      let that = this
      this.$api.post('/proxy/instance/update', this.instaces.edit.form).then(function (response) {
        that.instaces.edit.visit = false
        that.$message({message: '修改成功', type: 'success'})
        that.refresh_instances()
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    instance_remove_from_pool (id) {
      let that = this
      this.$api.post('/proxy/instance/pool/remove', {'id': id}).then(function (response) {
        that.$message({message: '移除成功', type: 'success'})
        that.refresh_instances()
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    instance_remove (id) {
      let that = this
      this.$api.post('/proxy/instance/remove', {'id': id}).then(function (response) {
        that.$message({message: '移除成功', type: 'success'})
        that.refresh_instances()
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    instance_into_pool (id) {
      let that = this
      this.$api.post('/proxy/instance/pool/add', {'id': id}).then(function (response) {
        that.$message({message: '加入成功', type: 'success'})
        that.refresh_instances()
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    refresh_server () {
      let that = this
      this.$api.get("/proxy/server/info").then(function (response) {
        that.server.info.data = response.detail.data
        that.server.info.instances = response.detail.instances
      })
    },
    refresh_instances () {
      let that = this
      this.$api.get("/proxy/instances").then(function (response) {
        that.instanceData = response.detail
      })
    },
    start_server () {
      this.loading = true
      let that = this
      this.$api.post("/proxy/server/start").then(function (response) {
        that.refresh_server()
        that.loading = false
        that.$notify({
          title: '启动成功',
          message: '服务器启动成功',
          type: 'success'
        });
      }).then(function (response) {
        that.$notify({
          title: '启动失败',
          message: response.message,
          type: 'warning'
        });
        that.loading = false
      })
    },
    stop_server () {
      this.loading = true
      let that = this
      this.$api.post("/proxy/server/stop").then(function (response) {
        that.refresh_server()
        that.loading = false
        that.$notify({
          title: '停止成功',
          message: '服务器启动成功',
          type: 'success'
        });
      }).then(function (response) {
        that.$notify({
          title: '停止失败',
          message: response.message,
          type: 'warning'
        });
        that.loading = false
      })
    },
    start () {
      if (this.proxy.Status) {
        this.$message.warning('服务已经启动')
      } else{
        this.start_server()
      }
    },
    stop () {
      if (!this.proxy.Status) {
        this.$message.warning('服务已经停止')
      } else{
        this.stop_server()
      }
    },
    submit_create_role () {
      let that = this
      this.$api.post('/proxy/role/add', this.roles.create.form).then(function (response) {
        that.roles.create.visit = false
        that.$message({message: '修改成功', type: 'success'})
        that.refresh_roles()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    remove_role (item) {
      let that = this
      this.$api.post('/proxy/role/remove', {id: item.id}).then(function (response) {
        that.$message({message: '删除成功', type: 'success'})
        that.refresh_roles()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    refresh_roles () {
      let that = this
      this.$api.get("/proxy/roles").then(function (response) {
        that.rolesData = response.detail
      })
    },
    filter_status(value, row) {
      return row.status === value;
    },
    refresh_instances () {
      let that = this
      this.$api.get("/proxy/instances").then(function (response) {
        that.instanceData = response.detail
      })
    },
    role_change: function (status) {
      if (status) {
        this.roles.create.style = ""
      } else {
        this.roles.create.style = "display: none"
      }
    },
    get_instance: function (id) {
      name = "默认转发"
      if (this.instanceData != null) {
        this.instanceData.forEach(function (item) {
          if (item.id === id) {
            name = item.address
          }
        })
      }
      return name
    }
  
  },
  created: function () {
    this.refresh_instances()
    this.refresh_server()
    this.refresh_roles()
  },
  mounted: function () {}
};
</script>

<style>
.el-card__header {
  padding: 5px;
}
.el-card__body {
  padding: 5px;
}

.el-dialog__header {
  padding: 10px 10px 5px;
  border-bottom: 1px solid whitesmoke;
}

.el-dialog__headerbtn {
  top: 12px;
}
.el-dialog__body {
  padding: 15px 10px;
}
.el-dialog__footer {
  border-top: 1px solid whitesmoke;
  padding: 5px 10px 10px;
}

.el-divider--horizontal {
  margin: 3px 0;
}
</style>
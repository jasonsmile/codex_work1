<template>
  <main class="page" :class="{ 'page--home': currentView === 'home' }">
    <section v-if="currentView === 'home'" class="home-view">
      <img class="brand-logo" :src="techLogo" alt="信息管理平台标识" />
      <h1>信息管理平台</h1>
      <p class="tagline">A simple platform for information management.</p>
      <el-button type="primary" size="large" @click="openManagement">进入信息管理</el-button>
      <footer class="home-footer">
        <span>Powered by Codex</span>
        <span>Copyright 信息中心团队</span>
      </footer>
    </section>

    <section v-else class="workspace">
      <aside class="sidebar">
        <div class="sidebar-title">信息管理平台</div>
        <el-menu :default-active="activeMenu" class="side-menu" @select="handleMenuSelect">
          <el-menu-item index="drugs">药品信息管理</el-menu-item>
          <el-menu-item index="specimens">标本留存信息</el-menu-item>
          <el-menu-item index="about">关于我们</el-menu-item>
        </el-menu>
      </aside>

      <section class="content">
        <header class="page-header">
          <div>
            <h1>{{ pageTitle }}</h1>
            <p>{{ pageSubtitle }}</p>
          </div>
        </header>

        <template v-if="activeMenu === 'drugs'">
          <section class="panel">
            <h2>新增药品</h2>
            <el-form
              ref="formRef"
              :model="form"
              :rules="rules"
              label-width="96px"
              class="drug-form"
            >
              <el-row :gutter="18">
                <el-col :xs="24" :md="12">
                  <el-form-item label="药品名称" prop="name">
                    <el-input v-model.trim="form.name" placeholder="请输入药品名称" clearable />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="生产厂家" prop="manufacturer">
                    <el-input v-model.trim="form.manufacturer" placeholder="请输入生产厂家" clearable />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="批准文号" prop="approvalNumber">
                    <el-input v-model.trim="form.approvalNumber" placeholder="请输入批准文号" clearable />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="规格" prop="specification">
                    <el-input v-model.trim="form.specification" placeholder="例如 0.25g*24片" clearable />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="价格" prop="price">
                    <el-input-number
                      v-model="form.price"
                      :min="0.01"
                      :precision="2"
                      :step="0.1"
                      controls-position="right"
                      class="number-input"
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="库存数量" prop="stock">
                    <el-input-number
                      v-model="form.stock"
                      :min="1"
                      :step="1"
                      :precision="0"
                      controls-position="right"
                      class="number-input"
                    />
                  </el-form-item>
                </el-col>
              </el-row>

              <div class="form-actions">
                <el-button type="primary" :loading="saving" @click="submitForm">提交</el-button>
                <el-button @click="resetForm">清空</el-button>
              </div>
            </el-form>
          </section>

          <section class="panel">
            <div class="table-toolbar">
              <h2>药品列表</h2>
              <div class="search-box">
                <el-input
                  v-model.trim="keyword"
                  placeholder="按药品名称搜索"
                  clearable
                  @clear="fetchDrugs"
                  @keyup.enter="fetchDrugs"
                />
                <el-button type="primary" plain @click="fetchDrugs">搜索</el-button>
              </div>
            </div>

            <el-table :data="drugs" v-loading="loading" border stripe empty-text="暂无药品数据">
              <el-table-column prop="name" label="药品名称" min-width="150" />
              <el-table-column prop="manufacturer" label="生产厂家" min-width="160" />
              <el-table-column prop="approvalNumber" label="批准文号" min-width="150" />
              <el-table-column prop="specification" label="规格" min-width="120" />
              <el-table-column prop="price" label="价格" width="110">
                <template #default="{ row }">￥{{ Number(row.price).toFixed(2) }}</template>
              </el-table-column>
              <el-table-column prop="stock" label="库存数量" width="110" />
              <el-table-column prop="createdAt" label="录入时间" min-width="180">
                <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
              </el-table-column>
            </el-table>
          </section>
        </template>

        <template v-else-if="activeMenu === 'specimens'">
          <section class="panel">
            <h2>添加申请单</h2>
            <el-form
              ref="specimenFormRef"
              :model="specimenForm"
              :rules="specimenRules"
              label-width="128px"
              class="specimen-form"
            >
              <el-row :gutter="18">
                <el-col :xs="24" :md="8">
                  <el-form-item label="姓名" prop="name">
                    <el-input v-model.trim="specimenForm.name" placeholder="请输入姓名" clearable />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="8">
                  <el-form-item label="性别" prop="gender">
                    <el-select v-model="specimenForm.gender" placeholder="请选择性别" class="full-input">
                      <el-option label="男" value="男" />
                      <el-option label="女" value="女" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="8">
                  <el-form-item label="年龄" prop="age">
                    <el-input-number
                      v-model="specimenForm.age"
                      :min="1"
                      :max="120"
                      :precision="0"
                      controls-position="right"
                      class="number-input"
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="8">
                  <el-form-item label="ID号" prop="idNumber">
                    <el-input v-model.trim="specimenForm.idNumber" placeholder="请输入ID号" clearable />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="8">
                  <el-form-item label="送检标本类型" prop="sampleType">
                    <el-select v-model="specimenForm.sampleType" placeholder="请选择标本类型" class="full-input">
                      <el-option label="组织" value="组织" />
                      <el-option label="血浆" value="血浆" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="8">
                  <el-form-item label="PD-L1表达" prop="pdl1Expression">
                    <div class="percent-input">
                      <el-input-number
                        v-model="specimenForm.pdl1Expression"
                        :min="0"
                        :max="100"
                        :precision="0"
                        controls-position="right"
                        class="number-input"
                      />
                      <span>%</span>
                    </div>
                  </el-form-item>
                </el-col>
                <el-col :xs="24">
                  <el-form-item label="病理类型" prop="pathologyType">
                    <el-radio-group v-model="specimenForm.pathologyType">
                      <el-radio v-for="item in pathologyTypes" :key="item" :label="item" :value="item" />
                    </el-radio-group>
                  </el-form-item>
                </el-col>
                <el-col :xs="24">
                  <el-form-item label="分期" prop="stage">
                    <el-radio-group v-model="specimenForm.stage">
                      <el-radio-button v-for="item in stages" :key="item" :label="item" :value="item" />
                    </el-radio-group>
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="8">
                  <el-form-item label="送检医师" prop="doctor">
                    <el-input v-model.trim="specimenForm.doctor" placeholder="请输入送检医师" clearable />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="8">
                  <el-form-item label="送检日期" prop="inspectionDate">
                    <el-date-picker
                      v-model="specimenForm.inspectionDate"
                      type="date"
                      value-format="YYYY-MM-DD"
                      placeholder="请选择送检日期"
                      class="full-input"
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="8">
                  <el-form-item label="驱动基因突变">
                    <el-input
                      v-model.trim="specimenForm.driverGeneMutation"
                      placeholder="例如 EGFR 19del"
                      clearable
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="末次治疗">
                    <el-input
                      v-model.trim="specimenForm.lastTreatment"
                      type="textarea"
                      :rows="3"
                      placeholder="请输入末次治疗"
                    />
                  </el-form-item>
                </el-col>
                <el-col :xs="24" :md="12">
                  <el-form-item label="后续治疗方案">
                    <el-input
                      v-model.trim="specimenForm.followUpTreatment"
                      type="textarea"
                      :rows="3"
                      placeholder="请输入后续治疗方案"
                    />
                  </el-form-item>
                </el-col>
              </el-row>

              <div class="form-actions specimen-actions">
                <el-button type="primary" :loading="specimenSaving" @click="submitSpecimenForm">提交</el-button>
                <el-button @click="resetSpecimenForm">清空</el-button>
              </div>
            </el-form>
          </section>

          <section class="panel">
            <div class="table-toolbar">
              <h2>申请单列表</h2>
              <el-button type="primary" plain @click="fetchSpecimens">刷新</el-button>
            </div>
            <el-table
              :data="specimenApplications"
              v-loading="specimenLoading"
              border
              stripe
              empty-text="暂无申请单数据"
            >
              <el-table-column prop="name" label="姓名" width="100" />
              <el-table-column prop="gender" label="性别" width="80" />
              <el-table-column prop="age" label="年龄" width="80" />
              <el-table-column prop="idNumber" label="ID号" min-width="140" />
              <el-table-column prop="sampleType" label="标本类型" width="110" />
              <el-table-column prop="pathologyType" label="病理类型" min-width="170" />
              <el-table-column prop="pdl1Expression" label="PD-L1" width="100">
                <template #default="{ row }">{{ row.pdl1Expression }}%</template>
              </el-table-column>
              <el-table-column prop="stage" label="分期" width="90" />
              <el-table-column prop="doctor" label="送检医师" width="120" />
              <el-table-column prop="inspectionDate" label="送检日期" width="130" />
              <el-table-column prop="driverGeneMutation" label="驱动基因突变" min-width="150" />
              <el-table-column prop="createdAt" label="创建时间" min-width="180">
                <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
              </el-table-column>
            </el-table>
          </section>
        </template>

        <template v-else-if="activeMenu === 'about'">
          <section class="panel about-panel">
            <img class="about-logo" :src="techLogo" alt="信息管理平台标识" />
            <h2>信息管理平台</h2>
            <p>A simple platform for information management.</p>
            <div class="contact-line">联系方式：tel=7934</div>
          </section>
        </template>
      </section>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import techLogo from './assets/tech-logo.png'

const API_BASE = '/api'

const createInitialForm = () => ({
  name: '',
  manufacturer: '',
  approvalNumber: '',
  specification: '',
  price: 0.01,
  stock: 1
})

const createInitialSpecimenForm = () => ({
  name: '',
  gender: '',
  age: 1,
  idNumber: '',
  sampleType: '',
  pathologyType: '',
  pdl1Expression: 0,
  driverGeneMutation: '',
  stage: '',
  lastTreatment: '',
  followUpTreatment: '',
  doctor: '',
  inspectionDate: ''
})

const formRef = ref()
const specimenFormRef = ref()
const form = reactive(createInitialForm())
const specimenForm = reactive(createInitialSpecimenForm())
const drugs = ref([])
const specimenApplications = ref([])
const keyword = ref('')
const loading = ref(false)
const saving = ref(false)
const specimenLoading = ref(false)
const specimenSaving = ref(false)
const currentView = ref('home')
const activeMenu = ref('drugs')

const pathologyTypes = ['腺癌', '鳞癌', '腺鳞癌', '大细胞神经内分泌癌', '小细胞肺癌', '其他']
const stages = ['I', 'II', 'III', 'IV']

const pageTitle = computed(() => {
  if (activeMenu.value === 'specimens') {
    return '标本留存信息'
  }
  if (activeMenu.value === 'about') {
    return '关于我们'
  }
  return '药品信息管理'
})
const pageSubtitle = computed(() =>
  ({
    specimens: '录入和查看标本留存申请单信息',
    about: '平台介绍和联系方式',
    drugs: '录入、查询和管理药品基础库存信息'
  })[activeMenu.value]
)

const menuRoutes = {
  drugs: '/drugs',
  specimens: '/specimens',
  about: '/about'
}

const routeMenus = {
  '/drugs': 'drugs',
  '/specimens': 'specimens',
  '/about': 'about'
}

const rules = {
  name: [{ required: true, message: '请输入药品名称', trigger: 'blur' }],
  price: [
    { required: true, message: '请输入价格', trigger: 'change' },
    { type: 'number', min: 0.01, message: '价格必须大于 0', trigger: 'change' }
  ],
  stock: [
    { required: true, message: '请输入库存数量', trigger: 'change' },
    { type: 'number', min: 1, message: '库存数量必须大于 0', trigger: 'change' }
  ]
}

const specimenRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  age: [{ type: 'number', min: 1, message: '年龄必须大于 0', trigger: 'change' }],
  idNumber: [{ required: true, message: '请输入ID号', trigger: 'blur' }],
  sampleType: [{ required: true, message: '请选择送检标本类型', trigger: 'change' }],
  pathologyType: [{ required: true, message: '请选择病理类型', trigger: 'change' }],
  pdl1Expression: [{ type: 'number', min: 0, max: 100, message: 'PD-L1表达必须在0到100之间', trigger: 'change' }],
  stage: [{ required: true, message: '请选择分期', trigger: 'change' }],
  doctor: [{ required: true, message: '请输入送检医师', trigger: 'blur' }],
  inspectionDate: [{ required: true, message: '请选择送检日期', trigger: 'change' }]
}

const fetchDrugs = async () => {
  loading.value = true
  try {
    const { data } = await axios.get(`${API_BASE}/drugs/get`, {
      params: { name: keyword.value || undefined }
    })
    drugs.value = data.data || []
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '查询失败')
  } finally {
    loading.value = false
  }
}

const fetchSpecimens = async () => {
  specimenLoading.value = true
  try {
    const { data } = await axios.get(`${API_BASE}/specimens/get`)
    specimenApplications.value = data.data || []
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '查询申请单失败')
  } finally {
    specimenLoading.value = false
  }
}

const openManagement = async () => {
  currentView.value = 'management'
  activeMenu.value = 'drugs'
  updateRoute('/drugs')
  await fetchDrugs()
}

const handleMenuSelect = async (index) => {
  activeMenu.value = index
  updateRoute(menuRoutes[index])
  if (index === 'drugs') {
    await fetchDrugs()
  }
  if (index === 'specimens') {
    await fetchSpecimens()
  }
}

const applyRoute = async () => {
  const path = window.location.pathname
  const matchedMenu = routeMenus[path]
  if (!matchedMenu) {
    currentView.value = 'home'
    activeMenu.value = 'drugs'
    return
  }

  currentView.value = 'management'
  activeMenu.value = matchedMenu
  if (matchedMenu === 'drugs') {
    await fetchDrugs()
  }
  if (matchedMenu === 'specimens') {
    await fetchSpecimens()
  }
}

const updateRoute = (path) => {
  if (!path || window.location.pathname === path) {
    return
  }
  window.history.pushState({}, '', path)
}

const submitForm = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) {
    return
  }

  saving.value = true
  try {
    await axios.post(`${API_BASE}/drugs/add`, {
      name: form.name,
      manufacturer: form.manufacturer,
      approvalNumber: form.approvalNumber,
      specification: form.specification,
      price: Number(form.price),
      stock: Number(form.stock)
    })
    ElMessage.success('保存成功')
    resetForm()
    await fetchDrugs()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const submitSpecimenForm = async () => {
  const valid = await specimenFormRef.value.validate().catch(() => false)
  if (!valid) {
    return
  }

  specimenSaving.value = true
  try {
    await axios.post(`${API_BASE}/specimens/add`, {
      ...specimenForm,
      age: Number(specimenForm.age),
      pdl1Expression: Number(specimenForm.pdl1Expression)
    })
    ElMessage.success('保存成功')
    resetSpecimenForm()
    await fetchSpecimens()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '保存申请单失败')
  } finally {
    specimenSaving.value = false
  }
}

const resetForm = () => {
  Object.assign(form, createInitialForm())
  formRef.value?.clearValidate()
}

const resetSpecimenForm = () => {
  Object.assign(specimenForm, createInitialSpecimenForm())
  specimenFormRef.value?.clearValidate()
}

const formatTime = (value) => {
  if (!value) {
    return '-'
  }
  return new Date(value).toLocaleString('zh-CN', { hour12: false })
}

onMounted(() => {
  applyRoute()
  window.addEventListener('popstate', applyRoute)
})

onUnmounted(() => {
  window.removeEventListener('popstate', applyRoute)
})

watch(
  [currentView, activeMenu],
  () => {
    document.title = currentView.value === 'home' ? '信息管理平台' : pageTitle.value
  },
  { immediate: true }
)
</script>

<template>
  <main class="page" :class="{ 'page--home': currentView === 'home', 'page--workspace': currentView !== 'home' }">
    <section v-if="currentView === 'home'" class="login-shell">
      <div class="login-mascot" :class="{ 'login-mascot--watching': mascotWatching }" aria-hidden="true">
        <div class="mascot-card mascot-card--purple">
          <span class="eye eye--left"></span>
          <span class="eye eye--right"></span>
        </div>
        <div class="mascot-card mascot-card--dark">
          <span class="eye eye--left"></span>
          <span class="eye eye--right"></span>
        </div>
        <div class="mascot-card mascot-card--orange">
          <span class="eye eye--left"></span>
          <span class="eye eye--right"></span>
        </div>
        <div class="mascot-card mascot-card--yellow">
          <span class="eye eye--left"></span>
          <span class="eye eye--right"></span>
          <span class="mascot-mouth"></span>
        </div>
      </div>

      <section class="home-view login-view">
        <img class="brand-logo" :src="techLogo" alt="信息管理平台标识" />
        <h1>信息管理平台</h1>
        <p class="tagline">A simple platform for information management.</p>

        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-position="top"
          class="login-form"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model.trim="loginForm.username"
              placeholder="请输入用户名"
              size="large"
              clearable
              @focus="loginFocused = true"
              @blur="loginFocused = false"
              @keyup.enter="submitLogin"
            />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input
              v-model.trim="loginForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              show-password
              @focus="loginFocused = true"
              @blur="loginFocused = false"
              @keyup.enter="submitLogin"
            />
          </el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loginLoading"
            @click="submitLogin"
          >
            登录
          </el-button>
        </el-form>
      </section>

      <footer class="home-footer">
        <span>Powered by Codex</span>
        <span>Copyright 信息中心团队</span>
      </footer>
    </section>

    <section v-else class="app-layout">
      <aside class="sidebar">
        <div class="sidebar-brand">
          <img class="sidebar-logo" :src="techLogo" alt="信息管理平台标识" />
          <span>信息管理平台</span>
        </div>
        <el-menu
          ref="sideMenuRef"
          :default-active="activeMenu"
          :default-openeds="defaultOpenedMenus"
          class="side-menu"
          @open="handleMenuOpen"
          @close="handleMenuClose"
          @select="handleMenuSelect"
        >
          <el-menu-item v-if="canViewMenu('home')" index="home">
            <el-icon><House /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item v-if="canViewMenu('drugs')" index="drugs">
            <el-icon><FirstAidKit /></el-icon>
            <span>药品信息管理</span>
          </el-menu-item>
          <el-menu-item v-if="canViewMenu('specimens')" index="specimens">
            <el-icon><Document /></el-icon>
            <span>标本留存信息</span>
          </el-menu-item>
          <el-menu-item v-if="canViewMenu('files')" index="files">
            <el-icon><FolderOpened /></el-icon>
            <span>媒体库（上传和下载）</span>
          </el-menu-item>
          <el-menu-item v-if="canViewMenu('about')" index="about">
            <el-icon><InfoFilled /></el-icon>
            <span>关于我们</span>
          </el-menu-item>
          <el-sub-menu v-if="canViewMenu('users')" index="superAdmin">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>超级管理员</span>
            </template>
            <el-menu-item index="users">
              <el-icon><UserFilled /></el-icon>
              <span>用户管理</span>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </aside>

      <section class="main-shell">
        <header class="top-bar">
          <div class="top-brand">
            <span v-if="breadcrumbParent" class="breadcrumb-root">{{ breadcrumbParent }}</span>
            <span v-if="breadcrumbParent" class="breadcrumb-separator">/</span>
            <span class="breadcrumb-current">{{ pageTitle }}</span>
          </div>
          <div class="top-user">
            <el-dropdown class="user-box" trigger="click" @command="handleUserCommand">
              <button class="user-trigger" type="button">
                <el-avatar :size="36" :src="currentUser?.headerImg">
                  {{ currentUser?.username?.slice(0, 1)?.toUpperCase() }}
                </el-avatar>
                <span class="user-name">{{ currentUser?.username }}</span>
                <span class="user-arrow" aria-hidden="true"></span>
              </button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">
                    <el-icon class="dropdown-icon"><SwitchButton /></el-icon>
                    <span>退出登录</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </header>

        <section class="content">

        <template v-if="activeMenu === 'home'">
          <section class="home-dashboard">
            <div class="home-dashboard-hero">
              <img class="home-dashboard-image" :src="homeHeroArt" alt="首页配图" />
              <div class="home-dashboard-copy">
                <p>你知道吗</p>
                <p>我顶着非常大智力障碍和健忘症在学习</p>
              </div>
            </div>
          </section>
        </template>

        <template v-else-if="activeMenu === 'drugs'">
          <section v-if="canCreateDrugs" class="panel action-panel">
            <el-button type="primary" @click="drugDrawerVisible = true">
              <el-icon><Plus /></el-icon>
              <span>新增药品</span>
            </el-button>
          </section>

          <el-drawer v-model="drugDrawerVisible" direction="rtl" size="520px" :show-close="false">
            <template #header>
              <div class="drawer-header">
                <h2>新增药品</h2>
                <div class="drawer-actions">
                  <el-button @click="cancelDrugDrawer">取消</el-button>
                  <el-button type="primary" :loading="saving" @click="submitForm">确定</el-button>
                </div>
              </div>
            </template>
            <el-form ref="formRef" :model="form" :rules="rules" label-width="96px" class="drug-form">
              <el-row :gutter="18">
                <el-col :span="24">
                  <el-form-item label="药品名称" prop="name">
                    <el-input v-model.trim="form.name" placeholder="请输入药品名称" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="生产厂家" prop="manufacturer">
                    <el-input v-model.trim="form.manufacturer" placeholder="请输入生产厂家" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="批准文号" prop="approvalNumber">
                    <el-input v-model.trim="form.approvalNumber" placeholder="请输入批准文号" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="规格" prop="specification">
                    <el-input v-model.trim="form.specification" placeholder="例如 0.25g*24粒" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
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
                <el-col :span="24">
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
            </el-form>
          </el-drawer>

          <section class="panel query-panel">
            <div class="search-box">
              <el-input
                v-model.trim="keyword"
                placeholder="按药品名称搜索"
                clearable
                @clear="fetchDrugs"
                @keyup.enter="fetchDrugs"
              />
              <el-button type="primary" @click="fetchDrugs">查询</el-button>
            </div>
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

            <el-table :data="pagedDrugs" v-loading="loading" border stripe empty-text="暂无药品数据">
              <el-table-column prop="name" label="药品名称" min-width="150" />
              <el-table-column prop="manufacturer" label="生产厂家" min-width="160" />
              <el-table-column prop="approvalNumber" label="批准文号" min-width="150" />
              <el-table-column prop="specification" label="规格" min-width="120" />
              <el-table-column prop="price" label="价格" width="110">
                <template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template>
              </el-table-column>
              <el-table-column prop="stock" label="库存数量" width="110" />
              <el-table-column prop="createdAt" label="录入时间" min-width="180">
                <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
              </el-table-column>
            </el-table>
            <div class="pagination-bar">
              <span>当前第 {{ drugPage }} 页 / 共 {{ drugTotalPages }} 页</span>
              <el-pagination
                layout="prev, pager, next"
                :current-page="drugPage"
                :page-size="pageSize"
                :total="drugs.length"
                @current-change="drugPage = $event"
              />
            </div>
          </section>
        </template>

        <template v-else-if="activeMenu === 'specimens'">
          <el-drawer v-model="specimenDrawerVisible" direction="rtl" size="640px" :show-close="false">
            <template #header>
              <div class="drawer-header">
                <h2>添加申请单</h2>
                <div class="drawer-actions">
                  <el-button @click="cancelSpecimenDrawer">取消</el-button>
                  <el-button type="primary" :loading="specimenSaving" @click="submitSpecimenForm">确定</el-button>
                </div>
              </div>
            </template>
            <el-form
              ref="specimenFormRef"
              :model="specimenForm"
              :rules="specimenRules"
              label-width="128px"
              class="specimen-form"
            >
              <el-row :gutter="18">
                <el-col :span="24">
                  <el-form-item label="姓名" prop="name">
                    <el-input v-model.trim="specimenForm.name" placeholder="请输入姓名" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="性别" prop="gender">
                    <el-select v-model="specimenForm.gender" placeholder="请选择性别" class="full-input">
                      <el-option label="男" value="男" />
                      <el-option label="女" value="女" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
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
                <el-col :span="24">
                  <el-form-item label="ID号" prop="idNumber">
                    <el-input v-model.trim="specimenForm.idNumber" placeholder="请输入ID号" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="送检标本类型" prop="sampleType">
                    <el-select v-model="specimenForm.sampleType" placeholder="请选择标本类型" class="full-input">
                      <el-option label="组织" value="组织" />
                      <el-option label="血浆" value="血浆" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
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
                <el-col :span="24">
                  <el-form-item label="送检医师" prop="doctor">
                    <el-input v-model.trim="specimenForm.doctor" placeholder="请输入送检医师" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
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
                <el-col :span="24">
                  <el-form-item label="驱动基因突变">
                    <el-input
                      v-model.trim="specimenForm.driverGeneMutation"
                      placeholder="例如 EGFR 19del"
                      clearable
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="末次治疗">
                    <el-input
                      v-model.trim="specimenForm.lastTreatment"
                      type="textarea"
                      :rows="3"
                      placeholder="请输入末次治疗"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
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
            </el-form>
          </el-drawer>

          <section class="panel query-panel specimen-query-panel">
            <el-form :model="specimenSearchForm" class="specimen-search" @submit.prevent>
              <label class="query-field">
                <span>姓名</span>
                <el-input
                  v-model.trim="specimenSearchForm.name"
                  placeholder="姓名"
                  clearable
                  @clear="fetchSpecimens"
                  @keyup.enter="fetchSpecimens"
                />
              </label>
              <label class="query-field">
                <span>ID号</span>
                <el-input
                  v-model.trim="specimenSearchForm.idNumber"
                  placeholder="ID号"
                  clearable
                  @clear="fetchSpecimens"
                  @keyup.enter="fetchSpecimens"
                />
              </label>
              <label class="query-field query-field--date">
                <span>送检日期</span>
                <el-date-picker
                  v-model="specimenSearchForm.inspectionDateRange"
                  type="daterange"
                  value-format="YYYY-MM-DD"
                  start-placeholder="开始日期"
                  end-placeholder="结束日期"
                  range-separator="至"
                  class="specimen-date-range"
                  @change="fetchSpecimens"
                />
              </label>
              <el-upload
                v-if="canCreateSpecimens"
                class="specimen-upload"
                accept=".xlsx,.xls"
                :auto-upload="false"
                :disabled="specimenImporting"
                :show-file-list="false"
                :on-change="confirmSpecimenExcelUpload"
              >
                <el-button :loading="specimenImporting">
                  <el-icon><Upload /></el-icon>
                  <span>批量上传</span>
                </el-button>
                <el-progress
                  v-if="specimenImporting || specimenImportProgress > 0"
                  class="specimen-upload-progress"
                  :percentage="specimenImportProgress"
                  :stroke-width="6"
                />
              </el-upload>
              <div class="query-actions">
                <el-button type="primary" @click="fetchSpecimens">
                  <el-icon><Search /></el-icon>
                  <span>查询</span>
                </el-button>
                <el-button @click="resetSpecimenSearch">
                  <el-icon><Refresh /></el-icon>
                  <span>重置</span>
                </el-button>
              </div>
            </el-form>
          </section>

          <section class="panel">
            <div class="table-toolbar">
              <el-button v-if="canCreateSpecimens" type="primary" @click="specimenDrawerVisible = true">
                <el-icon><Plus /></el-icon>
                <span>新增</span>
              </el-button>
            </div>
            <div class="table-toolbar specimen-toolbar">
              <el-form :model="specimenSearchForm" class="specimen-search" @submit.prevent>
                <label class="query-field">
                  <span>姓名</span>
                  <el-input
                    v-model.trim="specimenSearchForm.name"
                    placeholder="姓名"
                    clearable
                    @clear="fetchSpecimens"
                    @keyup.enter="fetchSpecimens"
                  />
                </label>
                <label class="query-field">
                  <span>ID号</span>
                  <el-input
                    v-model.trim="specimenSearchForm.idNumber"
                    placeholder="ID号"
                    clearable
                    @clear="fetchSpecimens"
                    @keyup.enter="fetchSpecimens"
                  />
                </label>
                <label class="query-field query-field--date">
                  <span>送检日期</span>
                  <el-date-picker
                    v-model="specimenSearchForm.inspectionDateRange"
                    type="daterange"
                    value-format="YYYY-MM-DD"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    range-separator="至"
                    class="specimen-date-range"
                    @change="fetchSpecimens"
                  />
                </label>
                <div class="query-actions">
                  <el-button type="primary" @click="fetchSpecimens">查询</el-button>
                  <el-button @click="resetSpecimenSearch">重置</el-button>
                </div>
              </el-form>
              <div class="specimen-add-row">
                <el-button v-if="canCreateSpecimens" type="primary" @click="specimenDrawerVisible = true">
                  <el-icon><Plus /></el-icon>
                  <span>新增加</span>
                </el-button>
              </div>
            </div>
            <el-table
              :data="pagedSpecimenApplications"
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
            <div class="pagination-bar">
              <span>当前第 {{ specimenPage }} 页 / 共 {{ specimenTotalPages }} 页</span>
              <el-pagination
                layout="prev, pager, next"
                :current-page="specimenPage"
                :page-size="pageSize"
                :total="specimenApplications.length"
                @current-change="specimenPage = $event"
              />
            </div>
          </section>
        </template>

        <template v-else-if="activeMenu === 'files'">
          <section class="panel action-panel">
            <div class="media-actions">
              <el-upload
                accept=".jpg,.jpeg,.png,.svg,.mp4,.txt,.sql,.xls,.xlsx,.doc,.docx"
                :auto-upload="false"
                :disabled="mediaUploading"
                :show-file-list="false"
                :on-change="handleMediaFileSelect"
              >
                <el-button type="primary" :loading="mediaUploading">
                  <el-icon><Upload /></el-icon>
                  <span>上传文件</span>
                </el-button>
              </el-upload>
              <el-button :loading="mediaLoading" @click="fetchMediaFiles">
                <el-icon><Refresh /></el-icon>
                <span>刷新</span>
              </el-button>
            </div>
          </section>

          <section class="panel">
            <el-table
              :key="`media-${mediaFiles.length}-${mediaPage}-${mediaListVersion}`"
              class="media-table"
              :data="pagedMediaFiles"
              v-loading="mediaLoading"
              header-cell-class-name="media-table-header"
              row-key="id"
              border
              stripe
              empty-text="暂无文件数据"
            >
              <el-table-column prop="displayName" label="文件名称" min-width="260" />
              <el-table-column prop="fileType" label="类型" width="130">
                <template #default="{ row }">
                  <el-tag :type="getMediaTypeTagType(row.fileType)">
                    {{ getMediaTypeLabel(row.fileType) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="标签" width="130">
                <template #default="{ row }">
                  <el-tag class="media-format-tag" type="info">{{ getMediaExtension(row) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="createdAt" label="日期" min-width="180">
                <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
              </el-table-column>
              <el-table-column label="操作" width="120" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="downloadMediaFile(row)">
                    <el-icon><Download /></el-icon>
                    <span>下载</span>
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <div class="pagination-bar">
              <span>当前第 {{ mediaPage }} 页 / 共 {{ mediaTotalPages }} 页</span>
              <el-pagination
                layout="prev, pager, next"
                :current-page="mediaPage"
                :page-size="pageSize"
                :total="mediaFiles.length"
                @current-change="mediaPage = $event"
              />
            </div>
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

        <template v-else-if="activeMenu === 'users'">
          <section class="panel">
            <div class="table-toolbar">
              <h2>用户管理</h2>
              <el-button type="primary" @click="userDrawerVisible = true">
                <el-icon><Plus /></el-icon>
                <span>新增用户</span>
              </el-button>
            </div>
          </section>

          <el-drawer v-model="userDrawerVisible" direction="rtl" size="520px" :show-close="false">
            <template #header>
              <div class="drawer-header">
                <h2>新增用户</h2>
                <div class="drawer-actions">
                  <el-button @click="cancelUserDrawer">取消</el-button>
                  <el-button type="primary" :loading="userSaving" @click="submitUserForm">确定</el-button>
                </div>
              </div>
            </template>
            <el-form
              ref="userFormRef"
              :model="userForm"
              :rules="userRules"
              label-width="96px"
              class="user-form"
            >
              <el-row :gutter="18">
                <el-col :span="24">
                  <el-form-item label="用户名" prop="username">
                    <el-input v-model.trim="userForm.username" placeholder="请输入用户名" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="密码" prop="password">
                    <el-input v-model.trim="userForm.password" type="password" placeholder="请输入密码" show-password />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="角色" prop="authorityId">
                    <el-select v-model="userForm.authorityId" class="full-input" placeholder="请选择角色">
                      <el-option label="管理员" :value="888" />
                      <el-option label="标本管理员" :value="777" />
                      <el-option label="只读用户" :value="999" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="手机号">
                    <el-input v-model.trim="userForm.phone" placeholder="请输入手机号" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="邮箱">
                    <el-input v-model.trim="userForm.email" placeholder="请输入邮箱" clearable />
                  </el-form-item>
                </el-col>
                <el-col :span="24">
                  <el-form-item label="状态" prop="enable">
                    <el-select v-model="userForm.enable" class="full-input" placeholder="请选择状态">
                      <el-option label="正常" :value="1" />
                      <el-option label="冻结" :value="2" />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>

            </el-form>
          </el-drawer>

          <section class="panel">
            <el-table :data="pagedUsers" v-loading="userLoading" border stripe empty-text="暂无用户数据">
              <el-table-column prop="username" label="用户名" min-width="140" />
              <el-table-column prop="authorityId" label="角色" min-width="130">
                <template #default="{ row }">{{ roleLabel(row.authorityId) }}</template>
              </el-table-column>
              <el-table-column prop="phone" label="手机号" min-width="140" />
              <el-table-column prop="email" label="邮箱" min-width="180" />
              <el-table-column prop="enable" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.enable === 1 ? 'success' : 'danger'">
                    {{ row.enable === 1 ? '正常' : '冻结' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="createdAt" label="创建时间" min-width="180">
                <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
              </el-table-column>
              <el-table-column label="操作" width="120" fixed="right">
                <template #default="{ row }">
                  <el-button type="danger" link @click="deleteUser(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div class="pagination-bar">
              <span>当前第 {{ userPage }} 页 / 共 {{ userTotalPages }} 页</span>
              <el-pagination
                layout="prev, pager, next"
                :current-page="userPage"
                :page-size="pageSize"
                :total="users.length"
                @current-change="userPage = $event"
              />
            </div>
          </section>
        </template>
      </section>
      </section>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Download, FirstAidKit, FolderOpened, House, InfoFilled, Plus, Refresh, Search, Setting, SwitchButton, Upload, UserFilled } from '@element-plus/icons-vue'
import techLogo from './assets/tech-logo.png'
import homeHeroArt from './assets/home-hero-art.png'

const API_BASE = '/api'
const USER_STORAGE_KEY = 'medical-info-current-user'
const LOGIN_ROUTE = '/login'

axios.interceptors.request.use((config) => {
  const user = readStoredUser()
  if (user?.token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${user.token}`
  }
  return config
})

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      handleAuthExpired()
    }
    return Promise.reject(error)
  }
)

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

const createInitialUserForm = () => ({
  username: '',
  password: '',
  authorityId: 999,
  phone: '',
  email: '',
  enable: 1
})

const loginFormRef = ref()
const formRef = ref()
const specimenFormRef = ref()
const userFormRef = ref()
const sideMenuRef = ref()
const loginForm = reactive({ username: '', password: '' })
const form = reactive(createInitialForm())
const specimenForm = reactive(createInitialSpecimenForm())
const specimenSearchForm = reactive({
  name: '',
  idNumber: '',
  inspectionDateRange: []
})
const userForm = reactive(createInitialUserForm())
const drugs = ref([])
const specimenApplications = ref([])
const users = ref([])
const mediaFiles = ref([])
const pageSize = 20
const drugPage = ref(1)
const specimenPage = ref(1)
const userPage = ref(1)
const mediaPage = ref(1)
const keyword = ref('')
const loading = ref(false)
const saving = ref(false)
const specimenLoading = ref(false)
const specimenSaving = ref(false)
const specimenImporting = ref(false)
const specimenImportProgress = ref(0)
const loginLoading = ref(false)
const loginFocused = ref(false)
const userLoading = ref(false)
const userSaving = ref(false)
const mediaLoading = ref(false)
const mediaUploading = ref(false)
const mediaListVersion = ref(0)
const drugDrawerVisible = ref(false)
const specimenDrawerVisible = ref(false)
const userDrawerVisible = ref(false)
const currentView = ref('home')
const activeMenu = ref('drugs')
const currentUser = ref(readStoredUser())
const defaultOpenedMenus = ref([])

const mascotWatching = computed(() => loginFocused.value || loginForm.username !== '' || loginForm.password !== '')

const roleMenus = {
  888: ['home', 'drugs', 'specimens', 'files', 'about', 'users'],
  777: ['home', 'specimens', 'files', 'about'],
  999: ['home', 'drugs', 'files', 'specimens', 'about']
}

const pathologyTypes = ['腺癌', '鳞癌', '腺鳞癌', '大细胞神经内分泌癌', '小细胞肺癌', '其他']
const stages = ['I', 'II', 'III', 'IV']

const pageTitle = computed(() => {
  if (activeMenu.value === 'home') {
    return '首页'
  }
  if (activeMenu.value === 'specimens') {
    return '标本留存信息'
  }
  if (activeMenu.value === 'about') {
    return '关于我们'
  }
  if (activeMenu.value === 'files') {
    return '媒体库（上传和下载）'
  }
  if (activeMenu.value === 'users') {
    return '用户管理'
  }
  return '药品信息管理'
})

const breadcrumbParent = computed(() => {
  if (activeMenu.value === 'users') {
    return '超级管理员'
  }
  return ''
})
const menuRoutes = {
  home: '/home',
  drugs: '/drugs',
  specimens: '/specimens',
  files: '/files',
  about: '/about',
  users: '/users'
}

const routeMenus = {
  '/home': 'home',
  '/drugs': 'drugs',
  '/specimens': 'specimens',
  '/files': 'files',
  '/about': 'about',
  '/users': 'users'
}

const allowedMenus = computed(() => {
  const authorityId = Number(currentUser.value?.authorityId)
  return roleMenus[authorityId] || []
})

const canCreateDrugs = computed(() => Number(currentUser.value?.authorityId) === 888)
const canCreateSpecimens = computed(() => [777, 888].includes(Number(currentUser.value?.authorityId)))

const canViewMenu = (menu) => allowedMenus.value.includes(menu)

const firstAllowedMenu = () => allowedMenus.value[0] || 'drugs'

const getTotalPages = (total) => Math.max(1, Math.ceil(total / pageSize))

const paginateRows = (rows, page) => {
  if (!Array.isArray(rows)) {
    return []
  }
  const start = (page - 1) * pageSize
  return rows.slice(start, start + pageSize)
}

const drugTotalPages = computed(() => getTotalPages(drugs.value.length))
const specimenTotalPages = computed(() => getTotalPages(specimenApplications.value.length))
const userTotalPages = computed(() => getTotalPages(users.value.length))
const mediaTotalPages = computed(() => getTotalPages(mediaFiles.value.length))

const pagedDrugs = computed(() => paginateRows(drugs.value, drugPage.value))
const pagedSpecimenApplications = computed(() => paginateRows(specimenApplications.value, specimenPage.value))
const pagedUsers = computed(() => paginateRows(users.value, userPage.value))
const pagedMediaFiles = computed(() => paginateRows(Array.isArray(mediaFiles.value) ? mediaFiles.value : [], mediaPage.value))

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
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

const userRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  authorityId: [{ required: true, message: '请选择角色', trigger: 'change' }],
  enable: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

function readStoredUser() {
  const raw = window.localStorage.getItem(USER_STORAGE_KEY)
  if (!raw) {
    return null
  }
  try {
    const user = JSON.parse(raw)
    if (!user?.token) {
      window.localStorage.removeItem(USER_STORAGE_KEY)
      return null
    }
    if (isTokenExpired(user.token)) {
      window.localStorage.removeItem(USER_STORAGE_KEY)
      return null
    }
    return user
  } catch {
    window.localStorage.removeItem(USER_STORAGE_KEY)
    return null
  }
}

function isTokenExpired(token) {
  const claims = parseTokenClaims(token)
  if (!claims?.expiresAt) {
    return true
  }
  return Number(claims.expiresAt) * 1000 <= Date.now()
}

function parseTokenClaims(token) {
  try {
    const payload = token.split('.')[0]
    if (!payload) {
      return null
    }
    const normalizedPayload = payload.replace(/-/g, '+').replace(/_/g, '/')
    const paddedPayload = normalizedPayload.padEnd(Math.ceil(normalizedPayload.length / 4) * 4, '=')
    return JSON.parse(window.atob(paddedPayload))
  } catch {
    return null
  }
}

const handleAuthExpired = () => {
  currentUser.value = null
  window.localStorage.removeItem(USER_STORAGE_KEY)
  currentView.value = 'home'
  activeMenu.value = 'drugs'
  updateRoute(LOGIN_ROUTE)
}

const getErrorMessage = (error, fallback) =>
  error.response?.data?.errorMessage || error.response?.data?.message || fallback

const fetchDrugs = async () => {
  loading.value = true
  try {
    const { data } = await axios.get(`${API_BASE}/drugs/get`, {
      params: { name: keyword.value || undefined }
    })
    drugs.value = data.data || []
    drugPage.value = 1
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '查询失败'))
  } finally {
    loading.value = false
  }
}

const fetchSpecimens = async () => {
  specimenLoading.value = true
  try {
    const [inspectionDateStart, inspectionDateEnd] = specimenSearchForm.inspectionDateRange || []
    const { data } = await axios.get(`${API_BASE}/specimens/get`, {
      params: {
        name: specimenSearchForm.name || undefined,
        idNumber: specimenSearchForm.idNumber || undefined,
        inspectionDateStart: inspectionDateStart || undefined,
        inspectionDateEnd: inspectionDateEnd || undefined
      }
    })
    specimenApplications.value = data.data || []
    specimenPage.value = 1
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '查询申请单失败'))
  } finally {
    specimenLoading.value = false
  }
}

const validateSpecimenExcelFile = (file) => {
  const fileName = file.name.toLowerCase()
  const validType = fileName.endsWith('.xlsx') || fileName.endsWith('.xls')
  if (!validType) {
    ElMessage.error('仅支持上传 .xlsx、.xls 格式文件')
    return false
  }
  if (file.size > 10 * 1024 * 1024) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }
  return true
}

const confirmSpecimenExcelUpload = async (uploadFile) => {
  const file = uploadFile.raw
  if (!file || specimenImporting.value) {
    return
  }
  if (!validateSpecimenExcelFile(file)) {
    return
  }

  try {
    const preview = await previewSpecimenExcel(file)
    await ElMessageBox.confirm(
      `<div class="upload-confirm-content">
        <p>确认上传文件 ${escapeHtml(file.name)}</p>
        <p>导入 <strong class="upload-count-success">${preview.successCount || 0}</strong> 条</p>
        <p>跳过 <strong class="upload-count-danger">${preview.skippedCount || 0}</strong> 条</p>
      </div>`,
      '批量上传确认',
      {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
      dangerouslyUseHTMLString: true
      }
    )
    await uploadSpecimenExcel(file)
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(getErrorMessage(error, '批量上传失败'))
    }
  }
}

const escapeHtml = (value) =>
  String(value).replace(/[&<>"']/g, (char) => ({
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;'
  })[char])

const previewSpecimenExcel = async (file) => {
  const formData = new FormData()
  formData.append('file', file)
  const { data } = await axios.post(`${API_BASE}/specimens/import/preview`, formData)
  return data.data || {}
}

const uploadSpecimenExcel = async (file) => {
  const formData = new FormData()
  formData.append('file', file)
  specimenImporting.value = true
  specimenImportProgress.value = 0
  try {
    const { data } = await axios.post(`${API_BASE}/specimens/import`, formData, {
      onUploadProgress: (event) => {
        if (event.total) {
          specimenImportProgress.value = Math.min(99, Math.round((event.loaded * 100) / event.total))
        }
      }
    })
    specimenImportProgress.value = 100
    const result = data.data || {}
    ElMessage.success(`导入完成，成功 ${result.successCount || 0} 条，跳过 ${result.skippedCount || 0} 条`)
    if (result.errors?.length) {
      const detail = result.errors
        .slice(0, 10)
        .map((item) => `第 ${item.row} 行：${item.message}`)
        .join('\n')
      await ElMessageBox.alert(detail, '导入跳过数据', {
        confirmButtonText: '知道了',
        type: 'warning'
      })
    }
    await fetchSpecimens()
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '批量上传失败'))
  } finally {
    specimenImporting.value = false
    window.setTimeout(() => {
      if (!specimenImporting.value) {
        specimenImportProgress.value = 0
      }
    }, 800)
  }
}

const fetchUsers = async () => {
  userLoading.value = true
  try {
    const { data } = await axios.get(`${API_BASE}/users/get`)
    users.value = data.data || []
    userPage.value = 1
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '查询用户失败'))
  } finally {
    userLoading.value = false
  }
}

const fetchMediaFiles = async () => {
  mediaLoading.value = true
  try {
    const { data } = await axios.get(`${API_BASE}/fileUploadAndDownload/get`)
    const list = Array.isArray(data.data) ? data.data : []
    mediaFiles.value = list.map((item) => ({ ...item }))
    mediaListVersion.value += 1
    mediaPage.value = 1
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '查询文件列表失败'))
  } finally {
    mediaLoading.value = false
  }
}

const validateMediaFile = (file) => {
  const fileName = file.name.toLowerCase()
  const validType = ['.jpg', '.jpeg', '.png', '.svg', '.mp4', '.txt', '.sql', '.xls', '.xlsx', '.doc', '.docx'].some((suffix) =>
    fileName.endsWith(suffix)
  )
  if (!validType) {
    ElMessage.error('仅支持 jpg、png、svg、mp4、txt、sql、xls、xlsx、doc、docx 格式文件')
    return false
  }
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('文件不能超过 5MB')
    return false
  }
  return true
}

const handleMediaFileSelect = async (uploadFile) => {
  const file = uploadFile.raw
  if (!file || mediaUploading.value) {
    return
  }
  if (!validateMediaFile(file)) {
    return
  }

  mediaUploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('displayName', file.name.replace(/\.[^.]+$/, ''))
    await axios.post(`${API_BASE}/fileUploadAndDownload/upload`, formData)
    ElMessage.success('上传成功')
    await fetchMediaFiles()
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '上传失败'))
  } finally {
    mediaUploading.value = false
  }
}

const downloadMediaFile = async (row) => {
  if (!row.id) {
    ElMessage.error('文件记录不存在')
    return
  }

  try {
    const { data } = await axios.get(`${API_BASE}/fileUploadAndDownload/download/${row.id}`)
    const downloadUrl = data.data?.url
    if (!downloadUrl) {
      ElMessage.error('文件地址不存在')
      return
    }
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = row.fileName
    link.target = '_blank'
    document.body.appendChild(link)
    link.click()
    link.remove()
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '下载失败'))
  }
}

const openMediaFile = (row) => {
  if (!row.url) {
    ElMessage.error('文件地址不存在')
    return
  }
  const link = document.createElement('a')
  link.href = row.url
  link.target = '_blank'
  document.body.appendChild(link)
  link.click()
  link.remove()
}

const submitLogin = async () => {
  const valid = await loginFormRef.value.validate().catch(() => false)
  if (!valid) {
    return
  }

  loginLoading.value = true
  try {
    const { data } = await axios.post(`${API_BASE}/users/login`, {
      username: loginForm.username,
      password: loginForm.password
    })
    currentUser.value = data.data
    window.localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(data.data))
    loginForm.password = ''
    ElMessage.success('登录成功')
    await openManagement()
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '登录失败'))
  } finally {
    loginLoading.value = false
  }
}

const openManagement = async () => {
  if (!currentUser.value) {
    currentView.value = 'home'
    updateRoute(LOGIN_ROUTE)
    return
  }

  currentView.value = 'management'
  activeMenu.value = firstAllowedMenu()
  updateRoute(menuRoutes[activeMenu.value])
  await fetchActiveMenuData()
}

const handleMenuSelect = async (index) => {
  if (!canViewMenu(index)) {
    ElMessage.error('无权访问该菜单')
    return
  }
  activeMenu.value = index
  syncSuperAdminMenu(index === 'users')
  updateRoute(menuRoutes[index])
  await fetchActiveMenuData()
}

const handleMenuOpen = (index) => {
  if (index === 'superAdmin') {
    defaultOpenedMenus.value = ['superAdmin']
  }
}

const handleMenuClose = (index) => {
  if (index === 'superAdmin') {
    defaultOpenedMenus.value = []
  }
}

const applyRoute = async () => {
  const path = window.location.pathname
  if (path === '/' || path === LOGIN_ROUTE) {
    if (currentUser.value) {
      currentView.value = 'management'
      activeMenu.value = 'home'
      updateRoute(menuRoutes.home)
      await fetchActiveMenuData()
      return
    }
    currentView.value = 'home'
    activeMenu.value = 'drugs'
    if (path !== LOGIN_ROUTE) {
      updateRoute(LOGIN_ROUTE)
    }
    return
  }

  const matchedMenu = routeMenus[path]
  if (!matchedMenu || !currentUser.value) {
    currentView.value = 'home'
    activeMenu.value = 'drugs'
    if (path !== LOGIN_ROUTE) {
      updateRoute(LOGIN_ROUTE)
    }
    return
  }

  if (!canViewMenu(matchedMenu)) {
    const fallbackMenu = firstAllowedMenu()
    currentView.value = 'management'
    activeMenu.value = fallbackMenu
    syncSuperAdminMenu(fallbackMenu === 'users')
    updateRoute(menuRoutes[fallbackMenu])
    await fetchActiveMenuData()
    return
  }

  currentView.value = 'management'
  activeMenu.value = matchedMenu
  syncSuperAdminMenu(matchedMenu === 'users')
  await fetchActiveMenuData()
}

const syncSuperAdminMenu = (shouldOpen) => {
  defaultOpenedMenus.value = shouldOpen ? ['superAdmin'] : []
  if (shouldOpen) {
    sideMenuRef.value?.open?.('superAdmin')
  } else {
    sideMenuRef.value?.close?.('superAdmin')
  }
}

const fetchActiveMenuData = async () => {
  if (activeMenu.value === 'drugs') {
    await fetchDrugs()
  }
  if (activeMenu.value === 'specimens') {
    await fetchSpecimens()
  }
  if (activeMenu.value === 'files') {
    await fetchMediaFiles()
  }
  if (activeMenu.value === 'users') {
    await fetchUsers()
  }
}

const updateRoute = (path) => {
  if (!path || window.location.pathname === path) {
    return
  }
  window.history.pushState({}, '', path)
}

const logout = () => {
  currentUser.value = null
  window.localStorage.removeItem(USER_STORAGE_KEY)
  currentView.value = 'home'
  activeMenu.value = 'drugs'
  updateRoute(LOGIN_ROUTE)
}

const handleUserCommand = (command) => {
  if (command === 'logout') {
    logout()
  }
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
    drugDrawerVisible.value = false
    await fetchDrugs()
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '保存失败'))
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
    specimenDrawerVisible.value = false
    await fetchSpecimens()
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '保存申请单失败'))
  } finally {
    specimenSaving.value = false
  }
}

const submitUserForm = async () => {
  const valid = await userFormRef.value.validate().catch(() => false)
  if (!valid) {
    return
  }

  userSaving.value = true
  try {
    await axios.post(`${API_BASE}/users/add`, {
      username: userForm.username,
      password: userForm.password,
      authorityId: Number(userForm.authorityId),
      phone: userForm.phone,
      email: userForm.email,
      enable: Number(userForm.enable)
    })
    ElMessage.success('创建用户成功')
    resetUserForm()
    userDrawerVisible.value = false
    await fetchUsers()
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '创建用户失败'))
  } finally {
    userSaving.value = false
  }
}

const deleteUser = async (row) => {
  if (row.id === currentUser.value?.id) {
    ElMessage.error('不能删除当前登录用户')
    return
  }

  try {
    await ElMessageBox.confirm(`确认删除用户 ${row.username} 吗？`, '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await axios.post(`${API_BASE}/users/delete`, { id: row.id })
    ElMessage.success('删除成功')
    await fetchUsers()
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    ElMessage.error(getErrorMessage(error, '删除用户失败'))
  }
}

const roleLabel = (authorityId) =>
  ({
    888: '管理员',
    777: '标本管理员',
    999: '只读用户'
  })[Number(authorityId)] || '未知角色'

const resetForm = () => {
  Object.assign(form, createInitialForm())
  formRef.value?.clearValidate()
}

const cancelDrugDrawer = () => {
  resetForm()
  drugDrawerVisible.value = false
}

const resetSpecimenForm = () => {
  Object.assign(specimenForm, createInitialSpecimenForm())
  specimenFormRef.value?.clearValidate()
}

const resetSpecimenSearch = async () => {
  specimenSearchForm.name = ''
  specimenSearchForm.idNumber = ''
  specimenSearchForm.inspectionDateRange = []
  await fetchSpecimens()
}

const cancelSpecimenDrawer = () => {
  resetSpecimenForm()
  specimenDrawerVisible.value = false
}

const resetUserForm = () => {
  Object.assign(userForm, createInitialUserForm())
  userFormRef.value?.clearValidate()
}

const cancelUserDrawer = () => {
  resetUserForm()
  userDrawerVisible.value = false
}

const getMediaExtension = (row) => {
  const fileName = row?.fileName || row?.displayName || ''
  const matched = fileName.match(/\.([^.]+)$/)
  if (matched?.[1]) {
    return matched[1].toLowerCase()
  }
  const contentType = row?.contentType || ''
  const typeMatched = contentType.match(/\/([a-z0-9.+-]+)$/i)
  return typeMatched?.[1]?.toLowerCase() || '-'
}

const getMediaTypeLabel = (fileType) => {
  const labels = {
    image: '图片',
    video: '视频',
    text: '文本'
  }
  return labels[fileType] || '其他'
}

const getMediaTypeTagType = (fileType) => {
  const types = {
    image: 'success',
    video: 'warning',
    text: 'primary'
  }
  return types[fileType] || 'info'
}

const formatTime = (value) => {
  if (!value) {
    return '-'
  }
  return new Date(value).toLocaleString()
}

const handlePopState = () => {
  applyRoute()
}

watch([currentView, activeMenu], () => {
  document.title = currentView.value === 'home' ? '信息管理平台登录' : pageTitle.value
})

watch(drugTotalPages, (totalPages) => {
  if (drugPage.value > totalPages) {
    drugPage.value = totalPages
  }
})

watch(specimenTotalPages, (totalPages) => {
  if (specimenPage.value > totalPages) {
    specimenPage.value = totalPages
  }
})

watch(userTotalPages, (totalPages) => {
  if (userPage.value > totalPages) {
    userPage.value = totalPages
  }
})

watch(mediaTotalPages, (totalPages) => {
  if (mediaPage.value > totalPages) {
    mediaPage.value = totalPages
  }
})

onMounted(() => {
  applyRoute()
  window.addEventListener('popstate', handlePopState)
})

onUnmounted(() => {
  window.removeEventListener('popstate', handlePopState)
})
</script>

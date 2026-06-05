<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import XuanwuDialog from '@/components/ui/XuanwuDialog.vue'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import Pagination from '@/components/Pagination.vue'
import { Plus, Pencil, Trash2, Eye, EyeOff, Search, AlertTriangle, Terminal, Zap, ZapOff } from 'lucide-vue-next'
import TextOverflow from '@/components/TextOverflow.vue'
import { api, type EnvVar } from '@/api'
import { toast } from 'vue-sonner'
import { useSiteSettings } from '@/composables/useSiteSettings'
import { Switch } from '@/components/ui/switch'
import { format } from 'date-fns'
import { Badge } from '@/components/ui/badge'

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  try {
    return format(new Date(dateStr), 'yyyy-MM-dd HH:mm:ss')
  } catch {
    return dateStr
  }
}

const { pageSize } = useSiteSettings()

const envVars = ref<EnvVar[]>([])
const showDialog = ref(false)
const editingEnv = ref<Partial<EnvVar>>({})
const isEdit = ref(false)
const showValues = ref<Record<string, boolean>>({})
const showDeleteDialog = ref(false)
const deleteEnvId = ref<string | null>(null)
const associatedTasks = ref<any[]>([])
const isDeleting = ref(false)
const valueTextareaRef = ref<HTMLTextAreaElement | null>(null)
const lineNumbersRef = ref<HTMLDivElement | null>(null)
const lineMeasureRef = ref<HTMLDivElement | null>(null)
const visualLineNumbers = ref<string[]>(['1'])
let textareaResizeObserver: ResizeObserver | null = null

const filterName = ref('')
const currentPage = ref(1)
const total = ref(0)
const activeTab = ref<string>('normal')
let searchTimer: ReturnType<typeof setTimeout> | null = null

async function loadEnvVars() {
  try {
    const res = await api.env.list({ page: currentPage.value, page_size: pageSize.value, name: filterName.value || undefined, type: activeTab.value })
    envVars.value = res.data
    total.value = res.total
    // 初始化显示状态，根据数据库的 hidden 状态同步显示
    res.data.forEach(env => {
      showValues.value[env.id] = !env.hidden
    })
  } catch { toast.error('加载环境变量失败') }
}

watch(showDialog, (val) => {
  if (!val) loadEnvVars()
})

function handleSearch() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    loadEnvVars()
  }, 300)
}

watch(activeTab, () => {
  currentPage.value = 1
  loadEnvVars()
})

function handlePageChange(page: number) {
  currentPage.value = page
  loadEnvVars()
}

function openCreate() {
  editingEnv.value = { name: '', value: '', remark: '', type: 'normal', hidden: true, enabled: true }
  isEdit.value = false
  showDialog.value = true
  void updateVisualLineNumbers()
}

function openEdit(env: EnvVar) {
  editingEnv.value = { ...env }
  isEdit.value = true
  showDialog.value = true
  void updateVisualLineNumbers()
}

function syncValueLineNumbers() {
  if (!valueTextareaRef.value || !lineNumbersRef.value) return
  lineNumbersRef.value.scrollTop = valueTextareaRef.value.scrollTop
}

async function updateVisualLineNumbers() {
  await nextTick()

  const textarea = valueTextareaRef.value
  const measure = lineMeasureRef.value
  if (!textarea || !measure) return

  const style = window.getComputedStyle(textarea)
  const lineHeight = Number.parseFloat(style.lineHeight) || Number.parseFloat(style.fontSize) * 1.5 || 24
  const lines = String(editingEnv.value?.value ?? '').split('\n')

  measure.style.width = `${textarea.clientWidth}px`
  measure.innerHTML = ''

  const nextLineNumbers: string[] = []
  lines.forEach((line, index) => {
    const lineEl = document.createElement('div')
    lineEl.className = 'break-all whitespace-pre-wrap'
    lineEl.textContent = line || ' '
    measure.appendChild(lineEl)

    const visualRows = Math.max(1, Math.round(lineEl.getBoundingClientRect().height / lineHeight))
    nextLineNumbers.push(String(index + 1))
    for (let i = 1; i < visualRows; i += 1) {
      nextLineNumbers.push('\u00A0')
    }
  })

  visualLineNumbers.value = nextLineNumbers.length > 0 ? nextLineNumbers : ['1']
  syncValueLineNumbers()
}

async function saveEnv() {
  try {
    if (isEdit.value && editingEnv.value.id) {
      await api.env.update(editingEnv.value.id, editingEnv.value)
      toast.success('变量已更新')
    } else {
      await api.env.create(editingEnv.value)
      toast.success('变量已创建')
    }
    showDialog.value = false
    loadEnvVars()
  } catch { toast.error('保存失败') }
}

async function confirmDelete(id: string) {
  deleteEnvId.value = id
  try {
    const res = await api.env.tasks(id)
    associatedTasks.value = res || []
    showDeleteDialog.value = true
  } catch {
    toast.error('检查变量引用失败')
  }
}

async function deleteEnv(force = false) {
  if (!deleteEnvId.value) return
  isDeleting.value = true
  try {
    const res = await api.env.delete(deleteEnvId.value, force)
    if (res.code === 409) {
      associatedTasks.value = res.data || []
      isDeleting.value = false
      return
    }
    if (res.code !== 200) {
      toast.error(res.msg || '删除失败')
      isDeleting.value = false
      return
    }
    toast.success('变量已删除')
    loadEnvVars()
    showDeleteDialog.value = false
  } catch {
    toast.error('网络错误，删除失败')
  } finally {
    isDeleting.value = false
  }
}

watch(showDeleteDialog, (val) => {
  if (!val) {
    associatedTasks.value = []
    deleteEnvId.value = null
  }
})

watch(() => editingEnv.value.value, () => {
  void updateVisualLineNumbers()
})

watch(showDialog, async (val) => {
  if (val) {
    await updateVisualLineNumbers()
    if (valueTextareaRef.value) {
      textareaResizeObserver?.disconnect()
      textareaResizeObserver = new ResizeObserver(() => {
        void updateVisualLineNumbers()
      })
      textareaResizeObserver.observe(valueTextareaRef.value)
    }
  } else {
    textareaResizeObserver?.disconnect()
  }
})

function toggleShow(id: string) {
  showValues.value[id] = !showValues.value[id]
}

async function toggleEnabled(env: EnvVar) {
  try {
    await api.env.update(env.id, { ...env, enabled: !env.enabled })
    env.enabled = !env.enabled
    toast.success(env.enabled ? '变量已启用' : '变量已禁用')
  } catch {
    toast.error('操作失败')
  }
}

function maskValue(value: string) {
  return '•'.repeat(Math.min(value.length, 20))
}

const NOTIFY_ENV_KEYS = ['XWPKG_NOTIFY_TOKEN', 'XWPKG_NOTIFY_CHANNEL', 'XWPKG_NOTIFY_URL']
function isNotifyEnv(name: string) {
  return NOTIFY_ENV_KEYS.includes(name)
}

onMounted(() => {
  loadEnvVars()
})

onBeforeUnmount(() => {
  textareaResizeObserver?.disconnect()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="flex flex-col shrink-0">
        <h2 class="text-xl sm:text-2xl font-bold tracking-tight">环境变量</h2>
        <p class="text-muted-foreground text-xs mt-0.5 ml-0.5">
          管理脚本执行时的环境变量
        </p>
      </div>

      <div class="flex flex-row items-center flex-wrap gap-2 w-full md:w-auto md:ml-auto md:justify-end">
        <!-- 搜索与操作 -->
        <div class="flex flex-row items-center gap-2 w-full sm:flex-1 md:flex-none md:w-auto text-sm">
          <div class="relative flex-1 md:flex-none md:w-[200px] group">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
            <Input v-model="filterName" placeholder="搜索名称..." class="h-9 pl-9 w-full bg-muted/20 border-muted-foreground/10 focus:bg-background text-sm" @input="handleSearch" />
          </div>
          
          <Button variant="outline" class="h-9 px-3 shrink-0 shadow-sm" @click="openCreate">
            <Plus class="h-4 w-4 md:mr-2" /> <span class="hidden md:inline">新建变量</span>
          </Button>
        </div>
      </div>
    </div>

    <div class="rounded-lg border bg-card overflow-hidden">
      <!-- ========== 1. 大屏布局 (Large >= 1280px) ========== -->
      <div class="hidden xl:block">
        <!-- 表头 -->
        <div class="flex items-center gap-4 px-4 py-1.5 border-b bg-muted/20 text-xs text-muted-foreground font-medium">
          <span class="w-12 shrink-0 pl-1">序号</span>
          <span class="w-48 shrink-0">名称</span>
          <span class="flex-1 min-w-0">值 / 内容</span>
          <span class="w-48 shrink-0">备注说明</span>
          <span class="w-40 shrink-0">创建时间</span>
          <span class="w-8 shrink-0 text-center">状态</span>
          <span class="w-24 shrink-0 text-center">操作</span>
        </div>
        <!-- 列表 -->
        <div class="divide-y text-sm">
          <div v-if="envVars.length === 0" class="text-sm text-muted-foreground text-center py-12">
            暂无环境变量
          </div>
          <div v-for="(env, index) in envVars" :key="`large-${env.id}`"
            class="flex items-center gap-4 px-4 py-1.5 hover:bg-muted/30 transition-colors">
            <div class="w-12 shrink-0 pl-1 text-muted-foreground tabular-nums">#{{ total - (currentPage - 1) * pageSize - index }}</div>
            
            <div class="w-48 shrink-0 flex items-center gap-1.5 overflow-hidden">
              <code class="font-bold truncate text-[11px] bg-muted/60 px-2 py-0.5 rounded text-zinc-700 dark:text-zinc-200">{{ env.name }}</code>
              <Badge v-if="isNotifyEnv(env.name)" variant="secondary" class="text-[9px] h-3.5 px-1 rounded-sm uppercase font-bold tracking-tighter shrink-0 leading-none">内置</Badge>
            </div>

            <div class="flex-1 min-w-0 text-muted-foreground truncate text-xs px-1">
              <TextOverflow :text="showValues[env.id] ? env.value : maskValue(env.value)" title="变量值" />
            </div>

            <div class="w-48 shrink-0 text-muted-foreground truncate text-xs">
              <TextOverflow :text="env.remark || '-'" title="备注描述" />
            </div>

            <div class="w-40 shrink-0 text-muted-foreground tabular-nums text-[11px] opacity-70">
              {{ formatDate(env.created_at) }}
            </div>

            <div class="w-8 shrink-0 flex justify-center">
              <span class="cursor-pointer" @click="toggleEnabled(env)">
                <div v-if="env.enabled" class="h-6 w-6 rounded-md bg-green-500/10 flex items-center justify-center">
                  <Zap class="h-3 w-3 text-green-500 fill-green-500" />
                </div>
                <div v-else class="h-6 w-6 rounded-md bg-muted flex items-center justify-center">
                  <ZapOff class="h-3 w-3 text-muted-foreground" />
                </div>
              </span>
            </div>

            <div class="w-24 shrink-0 flex items-center justify-center gap-1">
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="toggleShow(env.id)" :title="showValues[env.id] ? '隐藏' : '显示'">
                <Eye v-if="!showValues[env.id]" class="h-3.5 w-3.5" />
                <EyeOff v-else class="h-3.5 w-3.5" />
              </Button>
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="openEdit(env)" title="编辑">
                <Pencil class="h-3.5 w-3.5" />
              </Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive" @click="confirmDelete(env.id)" title="删除">
                <Trash2 class="h-3.5 w-3.5" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      <!-- ========== 2. 中屏布局 (Small to Large) ========== -->
      <div class="hidden sm:block xl:hidden">
        <!-- 表头 -->
        <div class="flex items-center gap-4 px-4 py-1.5 border-b bg-muted/20 text-xs text-muted-foreground font-medium">
          <span class="w-12 shrink-0 pl-1">序号</span>
          <span class="w-48 shrink-0">名称</span>
          <span class="flex-1 min-w-0">值 / 内容</span>
          <span class="w-8 shrink-0 text-center">状态</span>
          <span class="w-24 shrink-0 text-center">操作</span>
        </div>
        <!-- 列表 -->
        <div class="divide-y text-sm">
          <div v-for="(env, index) in envVars" :key="`med-${env.id}`"
            class="flex items-center gap-4 px-4 py-2 hover:bg-muted/30 transition-colors">
            <div class="w-12 shrink-0 pl-1 text-muted-foreground tabular-nums text-xs">#{{ total - (currentPage - 1) * pageSize - index }}</div>
            
            <div class="w-48 shrink-0 flex items-center gap-1.5 overflow-hidden">
              <code class="font-bold truncate text-[11px] bg-muted/60 px-2 py-0.5 rounded text-zinc-700 dark:text-zinc-200">{{ env.name }}</code>
              <Badge v-if="isNotifyEnv(env.name)" variant="secondary" class="text-[9px] h-3.5 px-1 rounded-sm uppercase font-bold tracking-tighter shrink-0 leading-none">内置</Badge>
            </div>

            <div class="flex-1 min-w-0 text-muted-foreground truncate text-xs">
               <TextOverflow :text="showValues[env.id] ? env.value : maskValue(env.value)" />
            </div>

            <div class="w-8 shrink-0 flex justify-center">
              <span class="cursor-pointer" @click="toggleEnabled(env)">
                <div v-if="env.enabled" class="h-6 w-6 rounded-md bg-green-500/10 flex items-center justify-center">
                  <Zap class="h-3 w-3 text-green-500 fill-green-500" />
                </div>
                <div v-else class="h-6 w-6 rounded-md bg-muted flex items-center justify-center">
                  <ZapOff class="h-3 w-3 text-muted-foreground" />
                </div>
              </span>
            </div>

            <div class="w-24 shrink-0 flex items-center justify-center gap-1">
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="toggleShow(env.id)">
                <Eye v-if="!showValues[env.id]" class="h-3.5 w-3.5" />
                <EyeOff v-else class="h-3.5 w-3.5" />
              </Button>
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="openEdit(env)">
                <Pencil class="h-3.5 w-3.5" />
              </Button>
              <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive" @click="confirmDelete(env.id)">
                <Trash2 class="h-3.5 w-3.5" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      <!-- ========== 3. 小屏布局 (Small < 640px) ========== -->
      <div class="divide-y sm:hidden">
        <div v-if="envVars.length === 0" class="text-sm text-muted-foreground text-center py-12">
          暂无环境变量
        </div>
        <div v-for="(env, index) in envVars" :key="`small-${env.id}`" class="p-3 hover:bg-muted/50 transition-colors">
          <div class="flex items-start justify-between mb-3 border-b border-border/40 pb-2">
            <div class="flex items-center gap-2 flex-1 min-w-0 pr-2">
              <span class="text-xs text-muted-foreground tabular-nums flex-shrink-0">#{{ total - (currentPage - 1) * pageSize - index }}</span>
              <code class="font-bold text-xs bg-muted/60 px-2 py-0.5 rounded truncate text-zinc-700 dark:text-zinc-200">{{ env.name }}</code>
              <Badge v-if="isNotifyEnv(env.name)" variant="secondary" class="text-[8px] h-3.5 px-1 rounded-sm uppercase font-bold tracking-tighter leading-none shrink-0">内置</Badge>
            </div>
            <span @click="toggleEnabled(env)" class="cursor-pointer">
              <div v-if="env.enabled" class="h-6 w-6 rounded-md bg-green-500/10 flex items-center justify-center">
                <Zap class="h-3.5 w-3.5 text-green-500 fill-green-500" />
              </div>
              <div v-else class="h-6 w-6 rounded-md bg-muted flex items-center justify-center">
                <ZapOff class="h-3.5 w-3.5 text-muted-foreground" />
              </div>
            </span>
          </div>
          
          <!-- 详情信息 -->
          <div class="space-y-1.5 text-xs text-muted-foreground mb-3 px-1">
            <div class="flex items-start gap-3">
              <span class="w-10 shrink-0 font-medium mt-0.5 opacity-70">内容:</span>
              <div class="flex-1 min-w-0 text-foreground break-all line-clamp-2">
                <TextOverflow :text="showValues[env.id] ? env.value : maskValue(env.value)" />
              </div>
            </div>
            <div v-if="env.remark" class="flex items-start gap-3">
              <span class="w-10 shrink-0 font-medium mt-0.5 opacity-70">备注:</span>
              <span class="flex-1 text-[11px] line-clamp-1">{{ env.remark }}</span>
            </div>
          </div>

          <div class="grid grid-cols-3 items-center pt-2 mt-2 border-t border-border/40 -mx-1">
            <Button variant="ghost" class="h-9 px-0 text-xs gap-1.5 hover:bg-primary/5 rounded-none" @click="toggleShow(env.id)">
              <Eye v-if="!showValues[env.id]" class="h-3.5 w-3.5" />
              <EyeOff v-else class="h-3.5 w-3.5" />
              {{ showValues[env.id] ? '隐藏' : '显示' }}
            </Button>
            <Button variant="ghost" class="h-9 px-0 text-xs gap-1.5 hover:bg-primary/5 rounded-none border-l border-border/10" @click="openEdit(env)">
              <Pencil class="h-3.5 w-3.5" />编辑
            </Button>
            <Button variant="ghost" class="h-9 px-0 text-xs gap-1.5 hover:bg-destructive/5 text-destructive rounded-none border-l border-border/10" @click="confirmDelete(env.id)">
              <Trash2 class="h-3.5 w-3.5" />删除
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="mt-4">
      <Pagination :total="total" :page="currentPage" @update:page="handlePageChange" />
    </div>


    <Dialog v-model:open="showDialog">
      <DialogContent class="w-[calc(100vw-2rem)] max-w-md min-w-0">
        <DialogHeader>
          <DialogTitle>{{ isEdit ? '编辑变量' : '新建变量' }}</DialogTitle>
          <DialogDescription class="sr-only">编辑变量的名称、值、备注以及启用和隐藏状态。</DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-2 min-w-0">
          <div class="space-y-2 min-w-0">
            <Label class="text-sm">变量名</Label>
            <Input v-model="editingEnv.name" class="h-9 w-full min-w-0 text-sm" placeholder="MY_VAR" />
          </div>
          <div class="space-y-2 min-w-0">
            <Label>变量值</Label>
            <div class="relative flex min-w-0 overflow-hidden rounded-md border border-input bg-transparent shadow-xs focus-within:border-ring focus-within:ring-ring/50 focus-within:ring-[3px]">
              <div ref="lineNumbersRef" class="flex max-h-40 w-6 shrink-0 flex-col overflow-hidden border-r border-border bg-muted/30 py-2 text-right font-mono text-[10px] leading-6 text-muted-foreground">
                <span v-for="(line, index) in visualLineNumbers" :key="`${index}-${line}`" class="block h-6 px-1">{{ line }}</span>
              </div>
              <textarea
                ref="valueTextareaRef"
                v-model="editingEnv.value"
                rows="5"
                placeholder="输入变量内容..."
                class="max-h-40 min-h-16 w-full min-w-0 resize-none overflow-x-hidden bg-transparent pl-2 pr-3 py-2 font-mono placeholder:font-sans text-sm leading-6 break-all whitespace-pre-wrap outline-none"
                @scroll="syncValueLineNumbers"
              />
              <div aria-hidden="true" class="pointer-events-none absolute bottom-1.5 right-1.5 h-3.5 w-3.5 opacity-45">
                <span class="absolute bottom-0 right-0 h-px w-3 rotate-[-45deg] bg-border" />
                <span class="absolute bottom-1 right-0.5 h-px w-2 rotate-[-45deg] bg-border/80" />
                <span class="absolute bottom-2 right-1 h-px w-1 rotate-[-45deg] bg-border/60" />
              </div>
            </div>
            <div
              ref="lineMeasureRef"
              aria-hidden="true"
              class="pointer-events-none invisible fixed left-0 top-0 -z-10 min-h-16 break-all whitespace-pre-wrap px-2 py-2 font-mono text-sm leading-6"
            />
          </div>
          <div class="space-y-2 min-w-0">
            <Label>备注</Label>
            <Textarea v-model="editingEnv.remark" class="w-full min-w-0 resize-none break-all text-sm placeholder:font-sans" rows="3" placeholder="变量用途说明..." />
          </div>
          <div class="flex items-center justify-between space-x-2 pt-2">
            <Label class="text-sm font-medium">隐藏变量值</Label>
            <Switch v-model="editingEnv.hidden" />
          </div>
          <div class="flex items-center justify-between space-x-2 pt-2">
            <Label class="text-sm font-medium">启用变量</Label>
            <Switch v-model="editingEnv.enabled" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showDialog = false">取消</Button>
          <Button @click="saveEnv">保存</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <XuanwuDialog v-model:open="showDeleteDialog" :title="associatedTasks.length > 0 ? '风险删除确认' : '确认删除'">
      <div v-if="associatedTasks.length > 0" class="space-y-4">
        <div class="flex items-start gap-4 p-4 rounded-xl bg-destructive/5 border border-destructive/10">
          <AlertTriangle class="h-5 w-5 text-destructive shrink-0" />
          <div class="space-y-1">
            <p class="text-sm font-bold text-destructive">环境变量正在使用中</p>
            <p class="text-[13px] text-muted-foreground/80 leading-relaxed">
              该变量已被以下任务引用，直接删除可能导致任务运行失败。建议先移除引用或选择"强制删除"。
            </p>
          </div>
        </div>

        <div class="space-y-2">
          <p class="text-[11px] font-bold text-muted-foreground uppercase tracking-widest px-1">关联任务 ({{ associatedTasks.length }})</p>
          <div class="bg-muted/30 rounded-xl p-2 max-h-48 overflow-y-auto space-y-1.5 border border-border/40">
            <div v-for="task in associatedTasks" :key="task.id"
              class="text-xs flex items-center justify-between bg-card p-2.5 rounded-lg border border-border/50 hover:border-primary/30 transition-all">
              <div class="flex items-center gap-2.5 min-w-0">
                <Terminal class="h-3.5 w-3.5 text-primary/70" />
                <span class="font-medium truncate">{{ task.name }}</span>
              </div>
              <code class="text-[10px] text-muted-foreground/70 font-mono bg-muted/50 px-1.5 py-0.5 rounded">{{ task.id }}</code>
            </div>
          </div>
        </div>

        <div class="p-4 rounded-xl bg-muted/20 border border-border/10">
          <p class="text-xs text-muted-foreground leading-relaxed italic">
            提示：选择强制删除将自动解除以上任务对该变量的绑定并执行物理删除。
          </p>
        </div>
      </div>
      <p v-else class="text-[15px] leading-relaxed text-muted-foreground">确定要删除此环境变量吗？此操作无法撤销，请谨慎操作。</p>

      <template #footer>
        <Button variant="ghost" :disabled="isDeleting" @click="showDeleteDialog = false">取消</Button>
        <Button v-if="associatedTasks.length > 0" variant="destructive" class="shadow-lg shadow-destructive/20" @click="deleteEnv(true)" :disabled="isDeleting">
          {{ isDeleting ? '删除中...' : '确认强制删除' }}
        </Button>
        <Button v-else variant="destructive" class="shadow-lg shadow-destructive/20" @click="deleteEnv(false)" :disabled="isDeleting">
          {{ isDeleting ? '删除中...' : '确认删除' }}
        </Button>
      </template>
    </XuanwuDialog>
  </div>
</template>

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fileAPI, folderAPI } from '../api'

export const useFileStore = defineStore('file', () => {
  const files = ref([])
  const folders = ref([])
  const currentFolder = ref(null)
  const loading = ref(false)
  const breadcrumbs = ref([{ id: null, name: '首页' }])

  async function fetchFiles(folderId = null) {
    loading.value = true
    try {
      const [filesRes, foldersRes] = await Promise.all([
        fileAPI.list({ folder_id: folderId }),
        folderAPI.list({ parent_id: folderId })
      ])
      files.value = filesRes
      folders.value = foldersRes
      currentFolder.value = folderId
    } finally {
      loading.value = false
    }
  }

  async function createFolder(name, parentId = null) {
    const folder = await folderAPI.create({ name, parent_id: parentId })
    folders.value.push(folder)
    return folder
  }

  async function deleteFile(id) {
    await fileAPI.delete(id)
    files.value = files.value.filter(f => f.id !== id)
  }

  async function deleteFolder(id) {
    await folderAPI.delete(id)
    folders.value = folders.value.filter(f => f.id !== id)
  }

  async function renameFile(id, name) {
    const updated = await fileAPI.rename(id, name)
    const idx = files.value.findIndex(f => f.id === id)
    if (idx !== -1) files.value[idx] = updated
    return updated
  }

  async function renameFolder(id, name) {
    const updated = await folderAPI.rename(id, name)
    const idx = folders.value.findIndex(f => f.id === id)
    if (idx !== -1) folders.value[idx] = updated
    return updated
  }

  function navigateToFolder(folder) {
    breadcrumbs.value.push({ id: folder.id, name: folder.name })
    fetchFiles(folder.id)
  }

  function navigateToRoot() {
    breadcrumbs.value = [{ id: null, name: '首页' }]
    fetchFiles(null)
  }

  function navigateBack() {
    if (breadcrumbs.value.length > 1) {
      breadcrumbs.value.pop()
      const prev = breadcrumbs.value[breadcrumbs.value.length - 1]
      fetchFiles(prev.id)
    }
  }

  return {
    files, folders, currentFolder, loading, breadcrumbs,
    fetchFiles, createFolder, deleteFile, deleteFolder,
    renameFile, renameFolder, navigateToFolder, navigateToRoot, navigateBack
  }
})
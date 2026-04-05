import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SplitMerge from '../views/SplitMerge.vue'

vi.mock('naive-ui', () => ({
  NButton: { name: 'NButton', template: '<button><slot /></button>', props: ['type', 'loading', 'disabled'] },
  NInputNumber: { name: 'NInputNumber', template: '<input type="number" />', props: ['value', 'min', 'max', 'step'] }
}))

describe('SplitMerge.vue', () => {
  it('renders split tab by default', () => {
    const wrapper = mount(SplitMerge)
    expect(wrapper.vm.activeTab).toBe('split')
  })

  it('switches to merge tab', async () => {
    const wrapper = mount(SplitMerge)
    wrapper.vm.activeTab = 'merge'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.activeTab).toBe('merge')
  })

  it('switches back to split tab', async () => {
    const wrapper = mount(SplitMerge)
    wrapper.vm.activeTab = 'merge'
    await wrapper.vm.$nextTick()
    wrapper.vm.activeTab = 'split'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.activeTab).toBe('split')
  })

  it('has default chunk size of 2MB', () => {
    const wrapper = mount(SplitMerge)
    expect(wrapper.vm.chunkSize).toBe(2 * 1024 * 1024)
  })

  it('selects split file', async () => {
    const wrapper = mount(SplitMerge)
    const mockFile = new File(['test content'], 'test.txt', { type: 'text/plain' })
    wrapper.vm.selectedSplitFile = mockFile
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.selectedSplitFile.name).toBe('test.txt')
    expect(wrapper.vm.splitChunks).toBe(1)
  })

  it('calculates split chunks for large file', () => {
    const wrapper = mount(SplitMerge)
    wrapper.vm.selectedSplitFile = new File([new ArrayBuffer(5 * 1024 * 1024)], 'big.bin')
    expect(wrapper.vm.splitChunks).toBe(3) // 5MB / 2MB = 3 chunks
  })

  it('adds merge files', async () => {
    const wrapper = mount(SplitMerge)
    wrapper.vm.mergeFiles = [
      new File(['part1'], 'test.part001of2.txt'),
      new File(['part2'], 'test.part002of2.txt')
    ]
    expect(wrapper.vm.mergeFiles.length).toBe(2)
  })

  it('removes a merge file', () => {
    const wrapper = mount(SplitMerge)
    wrapper.vm.mergeFiles = [
      new File(['part1'], 'a.part001of2.txt'),
      new File(['part2'], 'a.part002of2.txt')
    ]
    wrapper.vm.removeMergeFile(0)
    expect(wrapper.vm.mergeFiles.length).toBe(1)
  })

  it('performs split operation', async () => {
    const wrapper = mount(SplitMerge)
    const content = 'x'.repeat(3 * 1024 * 1024) // 3MB
    wrapper.vm.selectedSplitFile = new File([content], 'bigfile.dat')
    wrapper.vm.chunkSize = 1 * 1024 * 1024 // 1MB chunks

    await wrapper.vm.doSplit()

    expect(wrapper.vm.splitResult.length).toBe(3)
    expect(wrapper.vm.splitResult[0].name).toMatch(/\.part001of3/)
  })

  it('performs merge operation', async () => {
    const wrapper = mount(SplitMerge)
    wrapper.vm.mergeFiles = [
      new File(['hello'], 'test.part001of2.txt'),
      new File(['world'], 'test.part002of2.txt')
    ]

    await wrapper.vm.doMerge()

    expect(wrapper.vm.mergeComplete).toBe(true)
    expect(wrapper.vm.mergedFileName).toBe('test.txt')
    expect(wrapper.vm.mergedBlob?.size).toBe(10)
  })

  it('uses default name when chunk filename cannot be parsed', async () => {
    const wrapper = mount(SplitMerge)
    wrapper.vm.mergeFiles = [
      new File(['data'], 'non-standard-name.dat')
    ]

    await wrapper.vm.doMerge()

    expect(wrapper.vm.mergedFileName).toBe('merged_file')
  })
})

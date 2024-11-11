<template>
  <div class="relative w-full max-w-2xl p-8 bg-white rounded shadow-md">
    <h2 class="mb-6 text-xl font-bold text-center">
      {{ pathname === '/' ? '上传文件' : '没找到漫画，试试“上传文件？' }}
    </h2>
    <!-- 文件选择与拖拽区域 -->
    <div
      class="text-center transition duration-300 border-4 border-gray-300 border-dashed rounded cursor-pointer hover:border-blue-500"
      @click="fileInput.click()"
      @dragenter.prevent="onDragEnter"
      @dragover.prevent="onDragOver"
      @dragleave.prevent="onDragLeave"
      @drop.prevent="onDrop"
      :class="{ 'border-blue-500': isDragOver }"
      ref="dropArea"
    >
      <p class="p-20 text-gray-500">将文件拖拽到此处或点击选择文件</p>
      <input type="file" multiple class="hidden" ref="fileInput" @change="onFileChange" />
    </div>
    <!-- 预览选中文件列表 -->
    <div class="mt-6" v-if="filesToUpload.length > 0">
      <h3 class="mb-4 text-lg font-semibold">选择的文件</h3>
      <ul class="divide-y divide-gray-200">
        <li
          v-for="(file, index) in filesToUpload"
          :key="file.name + file.size + file.lastModified"
          class="flex items-center justify-between py-2"
        >
          <div class="flex items-center">
            <img
              v-if="file.type.startsWith('image/')"
              :src="file.previewUrl"
              class="object-cover w-12 h-12 mr-4 rounded"
              @load="revokeObjectURL(file.previewUrl)"
            />
            <div>
              <p class="font-medium">{{ file.name }}</p>
              <p class="text-sm text-gray-500">{{ formatFileSize(file.size) }}</p>
            </div>
          </div>
          <button class="text-red-500 hover:text-red-700" @click="removeFile(file)">
            &#10005;
          </button>
        </li>
      </ul>
    </div>
    <!-- 进度条背景 显示与隐藏有过渡效果，不可点击。-->
    <div
      class="absolute top-0 left-0 h-full transition-opacity duration-500 bg-blue-500 pointer-events-none opacity-30"
      :style="{ width: uploadProgress + '%' }"
      v-show="uploading"
    ></div>
    <!-- 上传按钮 -->
    <button
      class="w-full py-2 mt-6 text-white bg-blue-500 rounded hover:bg-blue-600 disabled:opacity-50"
      :disabled="filesToUpload.length === 0 || uploading"
      @click="uploadFiles"
    >
      {{ uploading ? '上传文件 ' + uploadProgress.toFixed(2) + '%' : '上传文件' }}
    </button>
    <!-- 上传结果显示区域 -->
    <div class="mt-4" v-if="result.message">
      <p class="text-lg text-center text-green-500" v-if="result.success">{{ result.message }}</p>
      <p class="text-lg text-center text-red-500" v-else>{{ result.message }}</p>
      <ul class="text-gray-700 list-disc list-inside" v-if="result.files && result.files.length">
        <li v-for="file in result.files" :key="file">{{ file }}</li>
      </ul>
    </div>
  </div>
  <div class="flex-grow bg-gray-400 place-holder"></div>
</template>

<script setup>
import { ref } from 'vue';

const pathname = window.location.pathname;

const fileInput = ref(null);
const filesToUpload = ref([]);
const isDragOver = ref(false);
const uploading = ref(false);
const uploadProgress = ref(0);
const result = ref({ message: '', files: [], success: false });

function onFileChange(event) {
  const files = event.target.files;
  handleFiles(files);
  event.target.value = '';
}

function onDragEnter() {
  isDragOver.value = true;
}

function onDragOver() {
  isDragOver.value = true;
}

function onDragLeave() {
  isDragOver.value = false;
}

function onDrop(event) {
  isDragOver.value = false;
  const files = event.dataTransfer.files;
  handleFiles(files);
}

function handleFiles(files) {
  for (let file of files) {
    // Avoid duplicates
    if (
      !filesToUpload.value.some(
        (f) =>
          f.name === file.name &&
          f.size === file.size &&
          f.lastModified === file.lastModified
      )
    ) {
      // If it's an image, create a preview URL
      if (file.type.startsWith('image/')) {
        file.previewUrl = URL.createObjectURL(file);
      }
      filesToUpload.value.push(file);
    }
  }
}

function removeFile(file) {
  filesToUpload.value = filesToUpload.value.filter((f) => f !== file);
  if (file.previewUrl) {
    URL.revokeObjectURL(file.previewUrl);
  }
}

function formatFileSize(size) {
  if (size > 1024 * 1024 * 1024) {
    return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
  } else if (size > 1024 * 1024) {
    return (size / (1024 * 1024)).toFixed(2) + ' MB';
  } else {
    return (size / 1024).toFixed(2) + ' KB';
  }
}

function uploadFiles() {
  if (filesToUpload.value.length === 0) return;
  uploading.value = true;
  uploadProgress.value = 0;
  const formData = new FormData();
  filesToUpload.value.forEach((file) => {
    formData.append('files', file);
  });
  const xhr = new XMLHttpRequest();
  xhr.open('POST', '/api/upload', true);
  xhr.upload.addEventListener('progress', (e) => {
    if (e.lengthComputable) {
      uploadProgress.value = (e.loaded / e.total) * 100;
    }
  });
  xhr.onload = function () {
    if (xhr.status === 200) {
      const data = JSON.parse(xhr.responseText);
      result.value.message = data.message;
      result.value.files = data.files;
      result.value.success = true;
      // Revoke object URLs
      filesToUpload.value.forEach((file) => {
        if (file.previewUrl) {
          URL.revokeObjectURL(file.previewUrl);
        }
      });
      // Clear files
      filesToUpload.value = [];
      location.reload();
    } else {
      const errorData = JSON.parse(xhr.responseText);
      result.value.message = errorData.error || '上传失败';
      result.value.success = false;
    }
    uploading.value = false;
    uploadProgress.value = 0;
  };
  xhr.onerror = function () {
    result.value.message = '上传失败: 网络错误';
    result.value.success = false;
    uploading.value = false;
    uploadProgress.value = 0;
  };
  xhr.send(formData);
}
</script>

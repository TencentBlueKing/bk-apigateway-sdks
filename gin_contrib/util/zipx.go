/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipDirectory 将指定目录压缩为ZIP文件
func ZipDirectory(srcDir, dstZip string, includeExt ...string) error {
	// 转换排除后缀为统一格式（带点的小写）
	includeMap := make(map[string]struct{})
	for _, ext := range includeExt {
		ext = strings.ToLower(ext)
		if !strings.HasPrefix(ext, ".") && ext != "" {
			ext = "." + ext
		}
		includeMap[ext] = struct{}{}
	}

	// 创建目标ZIP文件
	zipFile, err := os.Create(dstZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建ZIP写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历目录
	err = filepath.Walk(srcDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 获取文件后缀并统一为小写带点格式
		fileExt := strings.ToLower(filepath.Ext(filePath))
		// 检查是否在指明的后缀列表中
		if _, included := includeMap[fileExt]; !included {
			return nil // 跳过该文件
		}

		// 获取相对路径，并转换为ZIP兼容格式
		relPath, err := filepath.Rel(srcDir, filePath)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath) // 确保使用斜杠

		// 跳过根目录自身
		if relPath == "." {
			return nil
		}

		// 创建ZIP文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		// 处理目录：添加斜杠后缀并使用Store压缩方式
		if info.IsDir() {
			header.Name += "/"
			header.Method = zip.Store
		} else {
			header.Method = zip.Deflate
		}

		// 写入文件头
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件，写入内容
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

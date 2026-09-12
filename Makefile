.PHONY: all test race bench trace escape clean help

# 默认帮助
help:
	@echo "Go 后端学习效能脚本命令清单:"
	@echo "  make test        - 运行全部单元测试"
	@echo "  make race        - 数据竞争检测 (Race Detector)"
	@echo "  make bench       - 运行基准测试并分析内存分配"
	@echo "  make escape      - 运行逃逸分析 (查看堆栈分配决策)"
	@echo "  make trace       - 生成并提示可视化跟踪 (go tool trace)"
	@echo "  make clean       - 清理测试产物与编译临时文件"

# 单元测试
test:
	go test -v ./...

# 数据竞争检测
race:
	go test -race -v ./...

# 性能基准测试
bench:
	go test -bench=. -benchmem -run=none ./...

# 逃逸分析 (查看底层优化与逃逸原因)
escape:
	go build -gcflags="-m -l" ./...

# 清理构建与测试产物
clean:
	rm -f *.out *.pprof *.test trace.out cpu.pprof mem.pprof

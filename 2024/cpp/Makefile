# Advent of Code 2024 - C++ solutions.
#
#   make -s run day=day7 < input.txt   build and run one day on stdin
#   make test                          build and run every dayN_test.cpp
#   make test day=day7                 build and run one day's tests
#   make list                          list the available days
#   make clean                         remove build/
CXX ?= clang++
CXXFLAGS ?= -O2 -Wall -Wextra -Wpedantic
override CXXFLAGS += -std=c++23
BUILD := build

DAYS := $(sort $(basename $(filter-out %_test.cpp,$(wildcard day*.cpp))))
TESTS := $(sort $(basename $(wildcard day*_test.cpp)))
SELECTED_TESTS := $(addprefix $(BUILD)/,$(if $(day),$(day)_test,$(TESTS)))

.PHONY: all run test list clean

all: $(addprefix $(BUILD)/,$(DAYS) $(TESTS))

$(BUILD)/%: %.cpp | $(BUILD)
	@$(CXX) $(CXXFLAGS) -MMD -MP -MF $@.d $< -o $@

$(BUILD):
	@mkdir -p $@

run:
	$(if $(day),,$(error usage: make -s run day=dayN < input.txt))
	@$(MAKE) --no-print-directory $(BUILD)/$(day)
	@./$(BUILD)/$(day)

test:
	$(if $(SELECTED_TESTS),,$(error no dayN_test.cpp files found))
	@$(MAKE) --no-print-directory $(SELECTED_TESTS)
	@status=0; for t in $(SELECTED_TESTS); do \
		if ./$$t; then echo "PASS $$t"; else echo "FAIL $$t"; status=1; fi; \
	done; exit $$status

list:
	@echo $(DAYS)

clean:
	rm -rf $(BUILD)

-include $(wildcard $(BUILD)/*.d)

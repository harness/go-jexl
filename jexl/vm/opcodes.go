// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

type Opcode byte

const (
	OpInvalid Opcode = iota
	OpPush
	OpInt
	OpPop
	OpStore
	OpLoadVar
	OpLoadConst
	OpLoadField
	OpLoadFast
	OpLoadMethod
	OpLoadFunc
	OpLoadEnv
	OpFetch
	OpFetchField
	OpMethod
	OpTrue
	OpFalse
	OpNil
	OpNegate
	OpNot
	OpEqual
	OpEqualInt
	OpEqualString
	OpJump
	OpJumpIfTrue
	OpJumpIfFalse
	OpJumpIfNil
	OpJumpIfNotNil
	OpJumpIfFalsy
	OpJumpIfEnd
	OpJumpBackward
	OpIn
	OpMatches
	OpMatchesConst
	OpInOrMatches
	OpStartsWith
	OpEndsWith
	OpInstanceOf
	OpLess
	OpMore
	OpLessOrEqual
	OpMoreOrEqual
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpModulo
	OpExponent
	OpRange
	OpSlice
	OpCall
	OpCall0
	OpCall1
	OpCall2
	OpCall3
	OpCallN
	OpCallFast
	OpCallSafe
	OpCallTyped
	OpArray
	OpMap
	OpDeref
	OpIncrementIndex
	OpDecrementIndex
	OpIncrementCount
	OpGetIndex
	OpGetCount
	OpGetLen
	OpGetAcc
	OpSetAcc
	OpSetIndex
	OpPointer
	OpThrow
	OpBegin
	OpAnd
	OpOr
	OpBitOr
	OpBitXor
	OpBitAnd
	OpBitNot
	OpShiftLeft
	OpShiftRight
	OpShiftRightU
	OpStrictEqual
	OpStrictNotEqual
	OpLambda
	OpCallLambda
	OpAssign
	OpCompoundAssign
	OpIncrement
	OpDecrement
	OpForEach
	OpBreak
	OpContinue
	OpReturn
	OpTry
	OpEmpty
	OpSize
	OpSet
	OpNew
	OpStoreMember        // store value into object field: stack=[obj, val], arg=const field name
	OpStoreIndex         // store value into collection at key: stack=[obj, key, val]
	OpCompoundStoreIndex // compound assign to collection: stack=[obj, key, rhs], arg=const CompoundAssignOp with Op string (VarIndex unused)
	OpCompoundStoreEnv   // compound assign to env variable: stack=[rhs], arg=const string key
	OpEnd                // This opcode must be at the end of this list.
)

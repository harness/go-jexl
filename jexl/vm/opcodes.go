// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

type Opcode byte

const (
	OpInvalid Opcode = iota
	OpPush
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
	OpCall
	OpCall0
	OpCall1
	OpCall2
	OpCall3
	OpCallN
	OpCallFast
	OpArray
	OpMap
	OpDeref
	OpIncrementIndex
	OpDecrementIndex
	OpThrow
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
	OpReturn
	OpTry
	OpEmpty
	OpSize
	OpSet
	OpNew
	OpStoreIndex         // store value into collection at key: stack=[obj, key, val]
	OpCompoundStoreIndex // compound assign to collection: stack=[obj, key, rhs], arg=const CompoundAssignOp with Op string (VarIndex unused)
	OpCompoundStoreEnv   // compound assign to env variable: stack=[rhs], arg=const string key
	OpEnd                // This opcode must be at the end of this list.
)

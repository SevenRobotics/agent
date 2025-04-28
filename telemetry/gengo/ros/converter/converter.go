package converter

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go_agent/iface"
	"go_agent/telemetry/cmd/channel"
	gengo_actionlib "go_agent/telemetry/gengo/ros/actionlib"
	gengo_actionlib_msgs "go_agent/telemetry/gengo/ros/actionlib_msgs"
	gengo_actionlib_tutorials "go_agent/telemetry/gengo/ros/actionlib_tutorials"
	gengo_audio_common_msgs "go_agent/telemetry/gengo/ros/audio_common_msgs"
	gengo_base_local_planner "go_agent/telemetry/gengo/ros/base_local_planner"
	gengo_bond "go_agent/telemetry/gengo/ros/bond"
	gengo_controller_manager_msgs "go_agent/telemetry/gengo/ros/controller_manager_msgs"
	gengo_diagnostic_msgs "go_agent/telemetry/gengo/ros/diagnostic_msgs"
	gengo_dynamic_reconfigure "go_agent/telemetry/gengo/ros/dynamic_reconfigure"
	gengo_geometry_msgs "go_agent/telemetry/gengo/ros/geometry_msgs"
	gengo_mpc_local_planner_msgs "go_agent/telemetry/gengo/ros/mpc_local_planner_msgs"
	gengo_nav_msgs "go_agent/telemetry/gengo/ros/nav_msgs"
	gengo_plotjuggler_msgs "go_agent/telemetry/gengo/ros/plotjuggler_msgs"
	gengo_realsense2_camera "go_agent/telemetry/gengo/ros/realsense2_camera"
	gengo_roscpp "go_agent/telemetry/gengo/ros/roscpp"
	gengo_rosgraph_msgs "go_agent/telemetry/gengo/ros/rosgraph_msgs"
	gengo_rospy_tutorials "go_agent/telemetry/gengo/ros/rospy_tutorials"
	gengo_sensor_msgs "go_agent/telemetry/gengo/ros/sensor_msgs"
	gengo_smach_msgs "go_agent/telemetry/gengo/ros/smach_msgs"
	gengo_sound_play "go_agent/telemetry/gengo/ros/sound_play"
	gengo_std_msgs "go_agent/telemetry/gengo/ros/std_msgs"
	gengo_theora_image_transport "go_agent/telemetry/gengo/ros/theora_image_transport"
	gengo_turtle_actionlib "go_agent/telemetry/gengo/ros/turtle_actionlib"
	gengo_uuid_msgs "go_agent/telemetry/gengo/ros/uuid_msgs"
	proto_actionlib "go_agent/telemetry/genproto/ros/actionlib"
	proto_actionlib_msgs "go_agent/telemetry/genproto/ros/actionlib_msgs"
	proto_actionlib_tutorials "go_agent/telemetry/genproto/ros/actionlib_tutorials"
	proto_audio_common_msgs "go_agent/telemetry/genproto/ros/audio_common_msgs"
	proto_base_local_planner "go_agent/telemetry/genproto/ros/base_local_planner"
	proto_bond "go_agent/telemetry/genproto/ros/bond"
	proto_controller_manager_msgs "go_agent/telemetry/genproto/ros/controller_manager_msgs"
	proto_diagnostic_msgs "go_agent/telemetry/genproto/ros/diagnostic_msgs"
	proto_dynamic_reconfigure "go_agent/telemetry/genproto/ros/dynamic_reconfigure"
	proto_geometry_msgs "go_agent/telemetry/genproto/ros/geometry_msgs"
	proto_mpc_local_planner_msgs "go_agent/telemetry/genproto/ros/mpc_local_planner_msgs"
	proto_nav_msgs "go_agent/telemetry/genproto/ros/nav_msgs"
	proto_plotjuggler_msgs "go_agent/telemetry/genproto/ros/plotjuggler_msgs"
	proto_realsense2_camera "go_agent/telemetry/genproto/ros/realsense2_camera"
	proto_roscpp "go_agent/telemetry/genproto/ros/roscpp"
	proto_rosgraph_msgs "go_agent/telemetry/genproto/ros/rosgraph_msgs"
	proto_rospy_tutorials "go_agent/telemetry/genproto/ros/rospy_tutorials"
	proto_sensor_msgs "go_agent/telemetry/genproto/ros/sensor_msgs"
	proto_smach_msgs "go_agent/telemetry/genproto/ros/smach_msgs"
	proto_sound_play "go_agent/telemetry/genproto/ros/sound_play"
	proto_std_msgs "go_agent/telemetry/genproto/ros/std_msgs"
	proto_theora_image_transport "go_agent/telemetry/genproto/ros/theora_image_transport"
	proto_turtle_actionlib "go_agent/telemetry/genproto/ros/turtle_actionlib"
	proto_uuid_msgs "go_agent/telemetry/genproto/ros/uuid_msgs"
	"go_agent/utils"
)

func ConvertStateFeedback(rosMsg gengo_mpc_local_planner_msgs.StateFeedback) (proto_mpc_local_planner_msgs.StateFeedback, error) {

	ret := proto_mpc_local_planner_msgs.StateFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert StateFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float64, 0, len(rosMsg.State))
	for _, m := range rosMsg.State {
		t1 = append(t1, m)
	}

	ret.State = append(ret.State, t1...)

	return ret, nil
}

func ConvertOptimalControlResult(rosMsg gengo_mpc_local_planner_msgs.OptimalControlResult) (proto_mpc_local_planner_msgs.OptimalControlResult, error) {

	ret := proto_mpc_local_planner_msgs.OptimalControlResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert OptimalControlResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.DimStates = rosMsg.DimStates

	ret.DimControls = rosMsg.DimControls

	t3 := make([]float64, 0, len(rosMsg.TimeStates))
	for _, m := range rosMsg.TimeStates {
		t3 = append(t3, m)
	}

	ret.TimeStates = append(ret.TimeStates, t3...)

	t4 := make([]float64, 0, len(rosMsg.States))
	for _, m := range rosMsg.States {
		t4 = append(t4, m)
	}

	ret.States = append(ret.States, t4...)

	t5 := make([]float64, 0, len(rosMsg.TimeControls))
	for _, m := range rosMsg.TimeControls {
		t5 = append(t5, m)
	}

	ret.TimeControls = append(ret.TimeControls, t5...)

	t6 := make([]float64, 0, len(rosMsg.Controls))
	for _, m := range rosMsg.Controls {
		t6 = append(t6, m)
	}

	ret.Controls = append(ret.Controls, t6...)

	ret.OptimalSolutionFound = rosMsg.OptimalSolutionFound

	ret.CpuTime = rosMsg.CpuTime

	return ret, nil
}

func ConvertPosition2DInt(rosMsg gengo_base_local_planner.Position2DInt) (proto_base_local_planner.Position2DInt, error) {

	ret := proto_base_local_planner.Position2DInt{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Position2DInt, msg is nil")
	//}

	ret.X = rosMsg.X

	ret.Y = rosMsg.Y

	return ret, nil
}

func ConvertExtrinsics(rosMsg gengo_realsense2_camera.Extrinsics) (proto_realsense2_camera.Extrinsics, error) {

	ret := proto_realsense2_camera.Extrinsics{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Extrinsics, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float64, 0, len(rosMsg.Rotation))
	for _, m := range rosMsg.Rotation {
		t1 = append(t1, m)
	}

	ret.Rotation = append(ret.Rotation, t1...)

	t2 := make([]float64, 0, len(rosMsg.Translation))
	for _, m := range rosMsg.Translation {
		t2 = append(t2, m)
	}

	ret.Translation = append(ret.Translation, t2...)

	return ret, nil
}

func ConvertIMUInfo(rosMsg gengo_realsense2_camera.IMUInfo) (proto_realsense2_camera.IMUInfo, error) {

	ret := proto_realsense2_camera.IMUInfo{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert IMUInfo, msg is nil")
	//}

	ret.FrameId = rosMsg.FrameId

	t1 := make([]float64, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		t1 = append(t1, m)
	}

	ret.Data = append(ret.Data, t1...)

	t2 := make([]float64, 0, len(rosMsg.NoiseVariances))
	for _, m := range rosMsg.NoiseVariances {
		t2 = append(t2, m)
	}

	ret.NoiseVariances = append(ret.NoiseVariances, t2...)

	t3 := make([]float64, 0, len(rosMsg.BiasVariances))
	for _, m := range rosMsg.BiasVariances {
		t3 = append(t3, m)
	}

	ret.BiasVariances = append(ret.BiasVariances, t3...)

	return ret, nil
}

func ConvertMetadata(rosMsg gengo_realsense2_camera.Metadata) (proto_realsense2_camera.Metadata, error) {

	ret := proto_realsense2_camera.Metadata{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Metadata, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.JsonData = rosMsg.JsonData

	return ret, nil
}

func ConvertTestFeedback(rosMsg gengo_actionlib.TestFeedback) (proto_actionlib.TestFeedback, error) {

	ret := proto_actionlib.TestFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestFeedback, msg is nil")
	//}

	ret.Feedback = rosMsg.Feedback

	return ret, nil
}

func ConvertTestRequestAction(rosMsg gengo_actionlib.TestRequestAction) (proto_actionlib.TestRequestAction, error) {

	ret := proto_actionlib.TestRequestAction{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestRequestAction, msg is nil")
	//}

	var err error

	t0, err := ConvertTestRequestActionGoal(rosMsg.ActionGoal)
	ret.ActionGoal = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertTestRequestActionResult(rosMsg.ActionResult)
	ret.ActionResult = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTestRequestActionFeedback(rosMsg.ActionFeedback)
	ret.ActionFeedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTestRequestActionFeedback(rosMsg gengo_actionlib.TestRequestActionFeedback) (proto_actionlib.TestRequestActionFeedback, error) {

	ret := proto_actionlib.TestRequestActionFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestRequestActionFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTestRequestFeedback(rosMsg.Feedback)
	ret.Feedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTestRequestActionGoal(rosMsg gengo_actionlib.TestRequestActionGoal) (proto_actionlib.TestRequestActionGoal, error) {

	ret := proto_actionlib.TestRequestActionGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestRequestActionGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTestRequestGoal(rosMsg.Goal)
	ret.Goal = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTwoIntsActionResult(rosMsg gengo_actionlib.TwoIntsActionResult) (proto_actionlib.TwoIntsActionResult, error) {

	ret := proto_actionlib.TwoIntsActionResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwoIntsActionResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTwoIntsResult(rosMsg.Result)
	ret.Result = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTestRequestGoal(rosMsg gengo_actionlib.TestRequestGoal) (proto_actionlib.TestRequestGoal, error) {

	ret := proto_actionlib.TestRequestGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestRequestGoal, msg is nil")
	//}

	ret.TerminateStatus = rosMsg.TerminateStatus

	ret.IgnoreCancel = rosMsg.IgnoreCancel

	ret.ResultText = rosMsg.ResultText

	ret.TheResult = rosMsg.TheResult

	ret.IsSimpleClient = rosMsg.IsSimpleClient

	return ret, nil
}

func ConvertTestAction(rosMsg gengo_actionlib.TestAction) (proto_actionlib.TestAction, error) {

	ret := proto_actionlib.TestAction{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestAction, msg is nil")
	//}

	var err error

	t0, err := ConvertTestActionGoal(rosMsg.ActionGoal)
	ret.ActionGoal = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertTestActionResult(rosMsg.ActionResult)
	ret.ActionResult = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTestActionFeedback(rosMsg.ActionFeedback)
	ret.ActionFeedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTestActionFeedback(rosMsg gengo_actionlib.TestActionFeedback) (proto_actionlib.TestActionFeedback, error) {

	ret := proto_actionlib.TestActionFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestActionFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTestFeedback(rosMsg.Feedback)
	ret.Feedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTestRequestResult(rosMsg gengo_actionlib.TestRequestResult) (proto_actionlib.TestRequestResult, error) {

	ret := proto_actionlib.TestRequestResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestRequestResult, msg is nil")
	//}

	ret.TheResult = rosMsg.TheResult

	ret.IsSimpleServer = rosMsg.IsSimpleServer

	return ret, nil
}

func ConvertTestResult(rosMsg gengo_actionlib.TestResult) (proto_actionlib.TestResult, error) {

	ret := proto_actionlib.TestResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestResult, msg is nil")
	//}

	ret.Result = rosMsg.Result

	return ret, nil
}

func ConvertTwoIntsActionGoal(rosMsg gengo_actionlib.TwoIntsActionGoal) (proto_actionlib.TwoIntsActionGoal, error) {

	ret := proto_actionlib.TwoIntsActionGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwoIntsActionGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTwoIntsGoal(rosMsg.Goal)
	ret.Goal = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTwoIntsGoal(rosMsg gengo_actionlib.TwoIntsGoal) (proto_actionlib.TwoIntsGoal, error) {

	ret := proto_actionlib.TwoIntsGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwoIntsGoal, msg is nil")
	//}

	ret.A = rosMsg.A

	ret.B = rosMsg.B

	return ret, nil
}

func ConvertTwoIntsResult(rosMsg gengo_actionlib.TwoIntsResult) (proto_actionlib.TwoIntsResult, error) {

	ret := proto_actionlib.TwoIntsResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwoIntsResult, msg is nil")
	//}

	ret.Sum = rosMsg.Sum

	return ret, nil
}

func ConvertTwoIntsFeedback(rosMsg gengo_actionlib.TwoIntsFeedback) (proto_actionlib.TwoIntsFeedback, error) {

	ret := proto_actionlib.TwoIntsFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwoIntsFeedback, msg is nil")
	//}

	return ret, nil
}

func ConvertTestActionGoal(rosMsg gengo_actionlib.TestActionGoal) (proto_actionlib.TestActionGoal, error) {

	ret := proto_actionlib.TestActionGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestActionGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTestGoal(rosMsg.Goal)
	ret.Goal = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTestActionResult(rosMsg gengo_actionlib.TestActionResult) (proto_actionlib.TestActionResult, error) {

	ret := proto_actionlib.TestActionResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestActionResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTestResult(rosMsg.Result)
	ret.Result = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTestGoal(rosMsg gengo_actionlib.TestGoal) (proto_actionlib.TestGoal, error) {

	ret := proto_actionlib.TestGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestGoal, msg is nil")
	//}

	ret.Goal = rosMsg.Goal

	return ret, nil
}

func ConvertTestRequestActionResult(rosMsg gengo_actionlib.TestRequestActionResult) (proto_actionlib.TestRequestActionResult, error) {

	ret := proto_actionlib.TestRequestActionResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestRequestActionResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTestRequestResult(rosMsg.Result)
	ret.Result = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTestRequestFeedback(rosMsg gengo_actionlib.TestRequestFeedback) (proto_actionlib.TestRequestFeedback, error) {

	ret := proto_actionlib.TestRequestFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TestRequestFeedback, msg is nil")
	//}

	return ret, nil
}

func ConvertTwoIntsAction(rosMsg gengo_actionlib.TwoIntsAction) (proto_actionlib.TwoIntsAction, error) {

	ret := proto_actionlib.TwoIntsAction{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwoIntsAction, msg is nil")
	//}

	var err error

	t0, err := ConvertTwoIntsActionGoal(rosMsg.ActionGoal)
	ret.ActionGoal = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertTwoIntsActionResult(rosMsg.ActionResult)
	ret.ActionResult = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTwoIntsActionFeedback(rosMsg.ActionFeedback)
	ret.ActionFeedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTwoIntsActionFeedback(rosMsg gengo_actionlib.TwoIntsActionFeedback) (proto_actionlib.TwoIntsActionFeedback, error) {

	ret := proto_actionlib.TwoIntsActionFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwoIntsActionFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertTwoIntsFeedback(rosMsg.Feedback)
	ret.Feedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertGoalID(rosMsg gengo_actionlib_msgs.GoalID) (proto_actionlib_msgs.GoalID, error) {

	ret := proto_actionlib_msgs.GoalID{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GoalID, msg is nil")
	//}

	ret.Id = rosMsg.Id

	return ret, nil
}

func ConvertGoalStatus(rosMsg gengo_actionlib_msgs.GoalStatus) (proto_actionlib_msgs.GoalStatus, error) {

	ret := proto_actionlib_msgs.GoalStatus{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GoalStatus, msg is nil")
	//}

	var err error

	t0, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t0

	if err != nil {
		return ret, err
	}

	ret.Status = uint32(rosMsg.Status)

	ret.Text = rosMsg.Text

	return ret, nil
}

func ConvertGoalStatusArray(rosMsg gengo_actionlib_msgs.GoalStatusArray) (proto_actionlib_msgs.GoalStatusArray, error) {

	ret := proto_actionlib_msgs.GoalStatusArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GoalStatusArray, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.StatusList {
		c, err := ConvertGoalStatus(m)
		if err != nil {
			return ret, err
		}
		ret.StatusList = append(ret.StatusList, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertAveragingActionGoal(rosMsg gengo_actionlib_tutorials.AveragingActionGoal) (proto_actionlib_tutorials.AveragingActionGoal, error) {

	ret := proto_actionlib_tutorials.AveragingActionGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AveragingActionGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertAveragingGoal(rosMsg.Goal)
	ret.Goal = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertFibonacciActionResult(rosMsg gengo_actionlib_tutorials.FibonacciActionResult) (proto_actionlib_tutorials.FibonacciActionResult, error) {

	ret := proto_actionlib_tutorials.FibonacciActionResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert FibonacciActionResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertFibonacciResult(rosMsg.Result)
	ret.Result = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertFibonacciGoal(rosMsg gengo_actionlib_tutorials.FibonacciGoal) (proto_actionlib_tutorials.FibonacciGoal, error) {

	ret := proto_actionlib_tutorials.FibonacciGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert FibonacciGoal, msg is nil")
	//}

	ret.Order = rosMsg.Order

	return ret, nil
}

func ConvertAveragingAction(rosMsg gengo_actionlib_tutorials.AveragingAction) (proto_actionlib_tutorials.AveragingAction, error) {

	ret := proto_actionlib_tutorials.AveragingAction{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AveragingAction, msg is nil")
	//}

	var err error

	t0, err := ConvertAveragingActionGoal(rosMsg.ActionGoal)
	ret.ActionGoal = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertAveragingActionResult(rosMsg.ActionResult)
	ret.ActionResult = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertAveragingActionFeedback(rosMsg.ActionFeedback)
	ret.ActionFeedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertAveragingActionFeedback(rosMsg gengo_actionlib_tutorials.AveragingActionFeedback) (proto_actionlib_tutorials.AveragingActionFeedback, error) {

	ret := proto_actionlib_tutorials.AveragingActionFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AveragingActionFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertAveragingFeedback(rosMsg.Feedback)
	ret.Feedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertAveragingResult(rosMsg gengo_actionlib_tutorials.AveragingResult) (proto_actionlib_tutorials.AveragingResult, error) {

	ret := proto_actionlib_tutorials.AveragingResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AveragingResult, msg is nil")
	//}

	ret.Mean = rosMsg.Mean

	ret.StdDev = rosMsg.StdDev

	return ret, nil
}

func ConvertFibonacciFeedback(rosMsg gengo_actionlib_tutorials.FibonacciFeedback) (proto_actionlib_tutorials.FibonacciFeedback, error) {

	ret := proto_actionlib_tutorials.FibonacciFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert FibonacciFeedback, msg is nil")
	//}

	t0 := make([]int32, 0, len(rosMsg.Sequence))
	for _, m := range rosMsg.Sequence {
		t0 = append(t0, m)
	}

	ret.Sequence = append(ret.Sequence, t0...)

	return ret, nil
}

func ConvertFibonacciResult(rosMsg gengo_actionlib_tutorials.FibonacciResult) (proto_actionlib_tutorials.FibonacciResult, error) {

	ret := proto_actionlib_tutorials.FibonacciResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert FibonacciResult, msg is nil")
	//}

	t0 := make([]int32, 0, len(rosMsg.Sequence))
	for _, m := range rosMsg.Sequence {
		t0 = append(t0, m)
	}

	ret.Sequence = append(ret.Sequence, t0...)

	return ret, nil
}

func ConvertAveragingActionResult(rosMsg gengo_actionlib_tutorials.AveragingActionResult) (proto_actionlib_tutorials.AveragingActionResult, error) {

	ret := proto_actionlib_tutorials.AveragingActionResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AveragingActionResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertAveragingResult(rosMsg.Result)
	ret.Result = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertAveragingFeedback(rosMsg gengo_actionlib_tutorials.AveragingFeedback) (proto_actionlib_tutorials.AveragingFeedback, error) {

	ret := proto_actionlib_tutorials.AveragingFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AveragingFeedback, msg is nil")
	//}

	ret.Sample = rosMsg.Sample

	ret.Data = rosMsg.Data

	ret.Mean = rosMsg.Mean

	ret.StdDev = rosMsg.StdDev

	return ret, nil
}

func ConvertAveragingGoal(rosMsg gengo_actionlib_tutorials.AveragingGoal) (proto_actionlib_tutorials.AveragingGoal, error) {

	ret := proto_actionlib_tutorials.AveragingGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AveragingGoal, msg is nil")
	//}

	ret.Samples = rosMsg.Samples

	return ret, nil
}

func ConvertFibonacciActionFeedback(rosMsg gengo_actionlib_tutorials.FibonacciActionFeedback) (proto_actionlib_tutorials.FibonacciActionFeedback, error) {

	ret := proto_actionlib_tutorials.FibonacciActionFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert FibonacciActionFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertFibonacciFeedback(rosMsg.Feedback)
	ret.Feedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertFibonacciActionGoal(rosMsg gengo_actionlib_tutorials.FibonacciActionGoal) (proto_actionlib_tutorials.FibonacciActionGoal, error) {

	ret := proto_actionlib_tutorials.FibonacciActionGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert FibonacciActionGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertFibonacciGoal(rosMsg.Goal)
	ret.Goal = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertFibonacciAction(rosMsg gengo_actionlib_tutorials.FibonacciAction) (proto_actionlib_tutorials.FibonacciAction, error) {

	ret := proto_actionlib_tutorials.FibonacciAction{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert FibonacciAction, msg is nil")
	//}

	var err error

	t0, err := ConvertFibonacciActionGoal(rosMsg.ActionGoal)
	ret.ActionGoal = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertFibonacciActionResult(rosMsg.ActionResult)
	ret.ActionResult = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertFibonacciActionFeedback(rosMsg.ActionFeedback)
	ret.ActionFeedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertAudioData(rosMsg gengo_audio_common_msgs.AudioData) (proto_audio_common_msgs.AudioData, error) {

	ret := proto_audio_common_msgs.AudioData{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AudioData, msg is nil")
	//}

	t0 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := uint32(m)
		t0 = append(t0, a)
	}

	ret.Data = append(ret.Data, t0...)

	return ret, nil
}

func ConvertAudioDataStamped(rosMsg gengo_audio_common_msgs.AudioDataStamped) (proto_audio_common_msgs.AudioDataStamped, error) {

	ret := proto_audio_common_msgs.AudioDataStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AudioDataStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertAudioData(rosMsg.Audio)
	ret.Audio = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertAudioInfo(rosMsg gengo_audio_common_msgs.AudioInfo) (proto_audio_common_msgs.AudioInfo, error) {

	ret := proto_audio_common_msgs.AudioInfo{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AudioInfo, msg is nil")
	//}

	ret.Channels = uint32(rosMsg.Channels)

	ret.SampleRate = rosMsg.SampleRate

	ret.SampleFormat = rosMsg.SampleFormat

	ret.Bitrate = rosMsg.Bitrate

	ret.CodingFormat = rosMsg.CodingFormat

	return ret, nil
}

func ConvertConstants(rosMsg gengo_bond.Constants) (proto_bond.Constants, error) {

	ret := proto_bond.Constants{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Constants, msg is nil")
	//}

	return ret, nil
}

func ConvertStatus(rosMsg gengo_bond.Status) (proto_bond.Status, error) {

	ret := proto_bond.Status{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Status, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Id = rosMsg.Id

	ret.InstanceId = rosMsg.InstanceId

	ret.Active = rosMsg.Active

	ret.HeartbeatTimeout = rosMsg.HeartbeatTimeout

	ret.HeartbeatPeriod = rosMsg.HeartbeatPeriod

	return ret, nil
}

func ConvertControllersStatistics(rosMsg gengo_controller_manager_msgs.ControllersStatistics) (proto_controller_manager_msgs.ControllersStatistics, error) {

	ret := proto_controller_manager_msgs.ControllersStatistics{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ControllersStatistics, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Controller {
		c, err := ConvertControllerStatistics(m)
		if err != nil {
			return ret, err
		}
		ret.Controller = append(ret.Controller, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertHardwareInterfaceResources(rosMsg gengo_controller_manager_msgs.HardwareInterfaceResources) (proto_controller_manager_msgs.HardwareInterfaceResources, error) {

	ret := proto_controller_manager_msgs.HardwareInterfaceResources{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert HardwareInterfaceResources, msg is nil")
	//}

	ret.HardwareInterface = rosMsg.HardwareInterface

	t1 := make([]string, 0, len(rosMsg.Resources))
	for _, m := range rosMsg.Resources {
		t1 = append(t1, m)
	}

	ret.Resources = append(ret.Resources, t1...)

	return ret, nil
}

func ConvertControllerState(rosMsg gengo_controller_manager_msgs.ControllerState) (proto_controller_manager_msgs.ControllerState, error) {

	ret := proto_controller_manager_msgs.ControllerState{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ControllerState, msg is nil")
	//}

	var err error

	ret.Name = rosMsg.Name

	ret.State = rosMsg.State

	ret.Type = rosMsg.Type

	for _, m := range rosMsg.ClaimedResources {
		c, err := ConvertHardwareInterfaceResources(m)
		if err != nil {
			return ret, err
		}
		ret.ClaimedResources = append(ret.ClaimedResources, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertControllerStatistics(rosMsg gengo_controller_manager_msgs.ControllerStatistics) (proto_controller_manager_msgs.ControllerStatistics, error) {

	ret := proto_controller_manager_msgs.ControllerStatistics{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ControllerStatistics, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.Type = rosMsg.Type

	ret.Running = rosMsg.Running

	ret.NumControlLoopOverruns = rosMsg.NumControlLoopOverruns

	return ret, nil
}

func ConvertKeyValue(rosMsg gengo_diagnostic_msgs.KeyValue) (proto_diagnostic_msgs.KeyValue, error) {

	ret := proto_diagnostic_msgs.KeyValue{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert KeyValue, msg is nil")
	//}

	ret.Key = rosMsg.Key

	ret.Value = rosMsg.Value

	return ret, nil
}

func ConvertDiagnosticArray(rosMsg gengo_diagnostic_msgs.DiagnosticArray) (proto_diagnostic_msgs.DiagnosticArray, error) {

	ret := proto_diagnostic_msgs.DiagnosticArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert DiagnosticArray, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Status {
		c, err := ConvertDiagnosticStatus(m)
		if err != nil {
			return ret, err
		}
		ret.Status = append(ret.Status, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertDiagnosticStatus(rosMsg gengo_diagnostic_msgs.DiagnosticStatus) (proto_diagnostic_msgs.DiagnosticStatus, error) {

	ret := proto_diagnostic_msgs.DiagnosticStatus{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert DiagnosticStatus, msg is nil")
	//}

	var err error

	ret.Level = uint32(rosMsg.Level)

	ret.Name = rosMsg.Name

	ret.Message = rosMsg.Message

	ret.HardwareId = rosMsg.HardwareId

	for _, m := range rosMsg.Values {
		c, err := ConvertKeyValue(m)
		if err != nil {
			return ret, err
		}
		ret.Values = append(ret.Values, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertGroupState(rosMsg gengo_dynamic_reconfigure.GroupState) (proto_dynamic_reconfigure.GroupState, error) {

	ret := proto_dynamic_reconfigure.GroupState{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GroupState, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.State = rosMsg.State

	ret.Id = rosMsg.Id

	ret.Parent = rosMsg.Parent

	return ret, nil
}

func ConvertIntParameter(rosMsg gengo_dynamic_reconfigure.IntParameter) (proto_dynamic_reconfigure.IntParameter, error) {

	ret := proto_dynamic_reconfigure.IntParameter{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert IntParameter, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.Value = rosMsg.Value

	return ret, nil
}

func ConvertSensorLevels(rosMsg gengo_dynamic_reconfigure.SensorLevels) (proto_dynamic_reconfigure.SensorLevels, error) {

	ret := proto_dynamic_reconfigure.SensorLevels{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SensorLevels, msg is nil")
	//}

	return ret, nil
}

func ConvertBoolParameter(rosMsg gengo_dynamic_reconfigure.BoolParameter) (proto_dynamic_reconfigure.BoolParameter, error) {

	ret := proto_dynamic_reconfigure.BoolParameter{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert BoolParameter, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.Value = rosMsg.Value

	return ret, nil
}

func ConvertConfig(rosMsg gengo_dynamic_reconfigure.Config) (proto_dynamic_reconfigure.Config, error) {

	ret := proto_dynamic_reconfigure.Config{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Config, msg is nil")
	//}

	var err error

	for _, m := range rosMsg.Bools {
		c, err := ConvertBoolParameter(m)
		if err != nil {
			return ret, err
		}
		ret.Bools = append(ret.Bools, &c)
	}

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Ints {
		c, err := ConvertIntParameter(m)
		if err != nil {
			return ret, err
		}
		ret.Ints = append(ret.Ints, &c)
	}

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Strs {
		c, err := ConvertStrParameter(m)
		if err != nil {
			return ret, err
		}
		ret.Strs = append(ret.Strs, &c)
	}

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Doubles {
		c, err := ConvertDoubleParameter(m)
		if err != nil {
			return ret, err
		}
		ret.Doubles = append(ret.Doubles, &c)
	}

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Groups {
		c, err := ConvertGroupState(m)
		if err != nil {
			return ret, err
		}
		ret.Groups = append(ret.Groups, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertConfigDescription(rosMsg gengo_dynamic_reconfigure.ConfigDescription) (proto_dynamic_reconfigure.ConfigDescription, error) {

	ret := proto_dynamic_reconfigure.ConfigDescription{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ConfigDescription, msg is nil")
	//}

	var err error

	for _, m := range rosMsg.Groups {
		c, err := ConvertGroup(m)
		if err != nil {
			return ret, err
		}
		ret.Groups = append(ret.Groups, &c)
	}

	if err != nil {
		return ret, err
	}

	t1, err := ConvertConfig(rosMsg.Max)
	ret.Max = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertConfig(rosMsg.Min)
	ret.Min = &t2

	if err != nil {
		return ret, err
	}

	t3, err := ConvertConfig(rosMsg.Dflt)
	ret.Dflt = &t3

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertDoubleParameter(rosMsg gengo_dynamic_reconfigure.DoubleParameter) (proto_dynamic_reconfigure.DoubleParameter, error) {

	ret := proto_dynamic_reconfigure.DoubleParameter{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert DoubleParameter, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.Value = rosMsg.Value

	return ret, nil
}

func ConvertGroup(rosMsg gengo_dynamic_reconfigure.Group) (proto_dynamic_reconfigure.Group, error) {

	ret := proto_dynamic_reconfigure.Group{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Group, msg is nil")
	//}

	var err error

	ret.Name = rosMsg.Name

	ret.Type = rosMsg.Type

	for _, m := range rosMsg.Parameters {
		c, err := ConvertParamDescription(m)
		if err != nil {
			return ret, err
		}
		ret.Parameters = append(ret.Parameters, &c)
	}

	if err != nil {
		return ret, err
	}

	ret.Parent = rosMsg.Parent

	ret.Id = rosMsg.Id

	return ret, nil
}

func ConvertParamDescription(rosMsg gengo_dynamic_reconfigure.ParamDescription) (proto_dynamic_reconfigure.ParamDescription, error) {

	ret := proto_dynamic_reconfigure.ParamDescription{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ParamDescription, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.Type = rosMsg.Type

	ret.Level = rosMsg.Level

	ret.Description = rosMsg.Description

	ret.EditMethod = rosMsg.EditMethod

	return ret, nil
}

func ConvertStrParameter(rosMsg gengo_dynamic_reconfigure.StrParameter) (proto_dynamic_reconfigure.StrParameter, error) {

	ret := proto_dynamic_reconfigure.StrParameter{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert StrParameter, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.Value = rosMsg.Value

	return ret, nil
}

func ConvertAccelWithCovariance(rosMsg gengo_geometry_msgs.AccelWithCovariance) (proto_geometry_msgs.AccelWithCovariance, error) {

	ret := proto_geometry_msgs.AccelWithCovariance{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AccelWithCovariance, msg is nil")
	//}

	var err error

	t0, err := ConvertAccel(rosMsg.Accel)
	ret.Accel = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float64, 0, len(rosMsg.Covariance))
	for _, m := range rosMsg.Covariance {
		t1 = append(t1, m)
	}

	ret.Covariance = append(ret.Covariance, t1...)

	return ret, nil
}

func ConvertPolygon(rosMsg gengo_geometry_msgs.Polygon) (proto_geometry_msgs.Polygon, error) {

	ret := proto_geometry_msgs.Polygon{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Polygon, msg is nil")
	//}

	var err error

	for _, m := range rosMsg.Points {
		c, err := ConvertPoint32(m)
		if err != nil {
			return ret, err
		}
		ret.Points = append(ret.Points, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertPose(rosMsg gengo_geometry_msgs.Pose) (proto_geometry_msgs.Pose, error) {

	ret := proto_geometry_msgs.Pose{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Pose, msg is nil")
	//}

	var err error

	t0, err := ConvertPoint(rosMsg.Position)
	ret.Position = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertQuaternion(rosMsg.Orientation)
	ret.Orientation = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertPoseWithCovariance(rosMsg gengo_geometry_msgs.PoseWithCovariance) (proto_geometry_msgs.PoseWithCovariance, error) {

	ret := proto_geometry_msgs.PoseWithCovariance{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PoseWithCovariance, msg is nil")
	//}

	var err error

	t0, err := ConvertPose(rosMsg.Pose)
	ret.Pose = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float64, 0, len(rosMsg.Covariance))
	for _, m := range rosMsg.Covariance {
		t1 = append(t1, m)
	}

	ret.Covariance = append(ret.Covariance, t1...)

	return ret, nil
}

func ConvertTransform(rosMsg gengo_geometry_msgs.Transform) (proto_geometry_msgs.Transform, error) {

	ret := proto_geometry_msgs.Transform{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Transform, msg is nil")
	//}

	var err error

	t0, err := ConvertVector3(rosMsg.Translation)
	ret.Translation = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertQuaternion(rosMsg.Rotation)
	ret.Rotation = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTwistStamped(rosMsg gengo_geometry_msgs.TwistStamped) (proto_geometry_msgs.TwistStamped, error) {

	ret := proto_geometry_msgs.TwistStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwistStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertTwist(rosMsg.Twist)
	ret.Twist = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTwist(rosMsg gengo_geometry_msgs.Twist) (proto_geometry_msgs.Twist, error) {

	ret := proto_geometry_msgs.Twist{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Twist, msg is nil")
	//}

	var err error

	t0, err := ConvertVector3(rosMsg.Linear)
	ret.Linear = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertVector3(rosMsg.Angular)
	ret.Angular = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTwistWithCovariance(rosMsg gengo_geometry_msgs.TwistWithCovariance) (proto_geometry_msgs.TwistWithCovariance, error) {

	ret := proto_geometry_msgs.TwistWithCovariance{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwistWithCovariance, msg is nil")
	//}

	var err error

	t0, err := ConvertTwist(rosMsg.Twist)
	ret.Twist = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float64, 0, len(rosMsg.Covariance))
	for _, m := range rosMsg.Covariance {
		t1 = append(t1, m)
	}

	ret.Covariance = append(ret.Covariance, t1...)

	return ret, nil
}

func ConvertPoint(rosMsg gengo_geometry_msgs.Point) (proto_geometry_msgs.Point, error) {

	ret := proto_geometry_msgs.Point{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Point, msg is nil")
	//}

	ret.X = rosMsg.X

	ret.Y = rosMsg.Y

	ret.Z = rosMsg.Z

	return ret, nil
}

func ConvertPose2D(rosMsg gengo_geometry_msgs.Pose2D) (proto_geometry_msgs.Pose2D, error) {

	ret := proto_geometry_msgs.Pose2D{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Pose2D, msg is nil")
	//}

	ret.X = rosMsg.X

	ret.Y = rosMsg.Y

	ret.Theta = rosMsg.Theta

	return ret, nil
}

func ConvertPoseWithCovarianceStamped(rosMsg gengo_geometry_msgs.PoseWithCovarianceStamped) (proto_geometry_msgs.PoseWithCovarianceStamped, error) {

	ret := proto_geometry_msgs.PoseWithCovarianceStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PoseWithCovarianceStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertPoseWithCovariance(rosMsg.Pose)
	ret.Pose = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTransformStamped(rosMsg gengo_geometry_msgs.TransformStamped) (proto_geometry_msgs.TransformStamped, error) {

	ret := proto_geometry_msgs.TransformStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TransformStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.ChildFrameId = rosMsg.ChildFrameId

	t2, err := ConvertTransform(rosMsg.Transform)
	ret.Transform = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertVector3Stamped(rosMsg gengo_geometry_msgs.Vector3Stamped) (proto_geometry_msgs.Vector3Stamped, error) {

	ret := proto_geometry_msgs.Vector3Stamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Vector3Stamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertVector3(rosMsg.Vector)
	ret.Vector = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertWrench(rosMsg gengo_geometry_msgs.Wrench) (proto_geometry_msgs.Wrench, error) {

	ret := proto_geometry_msgs.Wrench{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Wrench, msg is nil")
	//}

	var err error

	t0, err := ConvertVector3(rosMsg.Force)
	ret.Force = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertVector3(rosMsg.Torque)
	ret.Torque = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertAccelWithCovarianceStamped(rosMsg gengo_geometry_msgs.AccelWithCovarianceStamped) (proto_geometry_msgs.AccelWithCovarianceStamped, error) {

	ret := proto_geometry_msgs.AccelWithCovarianceStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AccelWithCovarianceStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertAccelWithCovariance(rosMsg.Accel)
	ret.Accel = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertPoint32(rosMsg gengo_geometry_msgs.Point32) (proto_geometry_msgs.Point32, error) {

	ret := proto_geometry_msgs.Point32{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Point32, msg is nil")
	//}

	ret.X = rosMsg.X

	ret.Y = rosMsg.Y

	ret.Z = rosMsg.Z

	return ret, nil
}

func ConvertPolygonStamped(rosMsg gengo_geometry_msgs.PolygonStamped) (proto_geometry_msgs.PolygonStamped, error) {

	ret := proto_geometry_msgs.PolygonStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PolygonStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertPolygon(rosMsg.Polygon)
	ret.Polygon = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertQuaternionStamped(rosMsg gengo_geometry_msgs.QuaternionStamped) (proto_geometry_msgs.QuaternionStamped, error) {

	ret := proto_geometry_msgs.QuaternionStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert QuaternionStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertQuaternion(rosMsg.Quaternion)
	ret.Quaternion = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertQuaternion(rosMsg gengo_geometry_msgs.Quaternion) (proto_geometry_msgs.Quaternion, error) {

	ret := proto_geometry_msgs.Quaternion{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Quaternion, msg is nil")
	//}

	ret.X = rosMsg.X

	ret.Y = rosMsg.Y

	ret.Z = rosMsg.Z

	ret.W = rosMsg.W

	return ret, nil
}

func ConvertAccelStamped(rosMsg gengo_geometry_msgs.AccelStamped) (proto_geometry_msgs.AccelStamped, error) {

	ret := proto_geometry_msgs.AccelStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert AccelStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertAccel(rosMsg.Accel)
	ret.Accel = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertPointStamped(rosMsg gengo_geometry_msgs.PointStamped) (proto_geometry_msgs.PointStamped, error) {

	ret := proto_geometry_msgs.PointStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PointStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertPoint(rosMsg.Point)
	ret.Point = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertPoseArray(rosMsg gengo_geometry_msgs.PoseArray) (proto_geometry_msgs.PoseArray, error) {

	ret := proto_geometry_msgs.PoseArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PoseArray, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Poses {
		c, err := ConvertPose(m)
		if err != nil {
			return ret, err
		}
		ret.Poses = append(ret.Poses, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertPoseStamped(rosMsg gengo_geometry_msgs.PoseStamped) (proto_geometry_msgs.PoseStamped, error) {

	ret := proto_geometry_msgs.PoseStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PoseStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertPose(rosMsg.Pose)
	ret.Pose = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertInertia(rosMsg gengo_geometry_msgs.Inertia) (proto_geometry_msgs.Inertia, error) {

	ret := proto_geometry_msgs.Inertia{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Inertia, msg is nil")
	//}

	var err error

	ret.M = rosMsg.M

	t1, err := ConvertVector3(rosMsg.Com)
	ret.Com = &t1

	if err != nil {
		return ret, err
	}

	ret.Ixx = rosMsg.Ixx

	ret.Ixy = rosMsg.Ixy

	ret.Ixz = rosMsg.Ixz

	ret.Iyy = rosMsg.Iyy

	ret.Iyz = rosMsg.Iyz

	ret.Izz = rosMsg.Izz

	return ret, nil
}

func ConvertWrenchStamped(rosMsg gengo_geometry_msgs.WrenchStamped) (proto_geometry_msgs.WrenchStamped, error) {

	ret := proto_geometry_msgs.WrenchStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert WrenchStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertWrench(rosMsg.Wrench)
	ret.Wrench = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertAccel(rosMsg gengo_geometry_msgs.Accel) (proto_geometry_msgs.Accel, error) {

	ret := proto_geometry_msgs.Accel{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Accel, msg is nil")
	//}

	var err error

	t0, err := ConvertVector3(rosMsg.Linear)
	ret.Linear = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertVector3(rosMsg.Angular)
	ret.Angular = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertInertiaStamped(rosMsg gengo_geometry_msgs.InertiaStamped) (proto_geometry_msgs.InertiaStamped, error) {

	ret := proto_geometry_msgs.InertiaStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert InertiaStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertInertia(rosMsg.Inertia)
	ret.Inertia = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertTwistWithCovarianceStamped(rosMsg gengo_geometry_msgs.TwistWithCovarianceStamped) (proto_geometry_msgs.TwistWithCovarianceStamped, error) {

	ret := proto_geometry_msgs.TwistWithCovarianceStamped{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TwistWithCovarianceStamped, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertTwistWithCovariance(rosMsg.Twist)
	ret.Twist = &t1

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertVector3(rosMsg gengo_geometry_msgs.Vector3) (proto_geometry_msgs.Vector3, error) {

	ret := proto_geometry_msgs.Vector3{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Vector3, msg is nil")
	//}

	ret.X = rosMsg.X

	ret.Y = rosMsg.Y

	ret.Z = rosMsg.Z

	return ret, nil
}

func ConvertGetMapActionFeedback(rosMsg gengo_nav_msgs.GetMapActionFeedback) (proto_nav_msgs.GetMapActionFeedback, error) {

	ret := proto_nav_msgs.GetMapActionFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GetMapActionFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertGetMapFeedback(rosMsg.Feedback)
	ret.Feedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertGetMapActionGoal(rosMsg gengo_nav_msgs.GetMapActionGoal) (proto_nav_msgs.GetMapActionGoal, error) {

	ret := proto_nav_msgs.GetMapActionGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GetMapActionGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertGetMapGoal(rosMsg.Goal)
	ret.Goal = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertMapMetaData(rosMsg gengo_nav_msgs.MapMetaData) (proto_nav_msgs.MapMetaData, error) {

	ret := proto_nav_msgs.MapMetaData{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert MapMetaData, msg is nil")
	//}

	var err error

	ret.Resolution = rosMsg.Resolution

	ret.Width = rosMsg.Width

	ret.Height = rosMsg.Height

	t4, err := ConvertPose(rosMsg.Origin)
	ret.Origin = &t4

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertOdometry(rosMsg gengo_nav_msgs.Odometry) (proto_nav_msgs.Odometry, error) {

	ret := proto_nav_msgs.Odometry{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Odometry, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.ChildFrameId = rosMsg.ChildFrameId

	t2, err := ConvertPoseWithCovariance(rosMsg.Pose)
	ret.Pose = &t2

	if err != nil {
		return ret, err
	}

	t3, err := ConvertTwistWithCovariance(rosMsg.Twist)
	ret.Twist = &t3

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertGetMapAction(rosMsg gengo_nav_msgs.GetMapAction) (proto_nav_msgs.GetMapAction, error) {

	ret := proto_nav_msgs.GetMapAction{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GetMapAction, msg is nil")
	//}

	var err error

	t0, err := ConvertGetMapActionGoal(rosMsg.ActionGoal)
	ret.ActionGoal = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGetMapActionResult(rosMsg.ActionResult)
	ret.ActionResult = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertGetMapActionFeedback(rosMsg.ActionFeedback)
	ret.ActionFeedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertGetMapActionResult(rosMsg gengo_nav_msgs.GetMapActionResult) (proto_nav_msgs.GetMapActionResult, error) {

	ret := proto_nav_msgs.GetMapActionResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GetMapActionResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertGetMapResult(rosMsg.Result)
	ret.Result = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertGetMapFeedback(rosMsg gengo_nav_msgs.GetMapFeedback) (proto_nav_msgs.GetMapFeedback, error) {

	ret := proto_nav_msgs.GetMapFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GetMapFeedback, msg is nil")
	//}

	return ret, nil
}

func ConvertGetMapGoal(rosMsg gengo_nav_msgs.GetMapGoal) (proto_nav_msgs.GetMapGoal, error) {

	ret := proto_nav_msgs.GetMapGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GetMapGoal, msg is nil")
	//}

	return ret, nil
}

func ConvertGetMapResult(rosMsg gengo_nav_msgs.GetMapResult) (proto_nav_msgs.GetMapResult, error) {

	ret := proto_nav_msgs.GetMapResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GetMapResult, msg is nil")
	//}

	var err error

	t0, err := ConvertOccupancyGrid(rosMsg.Map)
	ret.Map = &t0

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertGridCells(rosMsg gengo_nav_msgs.GridCells) (proto_nav_msgs.GridCells, error) {

	ret := proto_nav_msgs.GridCells{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert GridCells, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.CellWidth = rosMsg.CellWidth

	ret.CellHeight = rosMsg.CellHeight

	for _, m := range rosMsg.Cells {
		c, err := ConvertPoint(m)
		if err != nil {
			return ret, err
		}
		ret.Cells = append(ret.Cells, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertOccupancyGrid(rosMsg gengo_nav_msgs.OccupancyGrid) (proto_nav_msgs.OccupancyGrid, error) {

	ret := proto_nav_msgs.OccupancyGrid{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert OccupancyGrid, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertMapMetaData(rosMsg.Info)
	ret.Info = &t1

	if err != nil {
		return ret, err
	}

	t2 := make([]int32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := int32(m)
		t2 = append(t2, a)
	}

	ret.Data = append(ret.Data, t2...)

	return ret, nil
}

func ConvertPath(rosMsg gengo_nav_msgs.Path) (proto_nav_msgs.Path, error) {

	ret := proto_nav_msgs.Path{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Path, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Poses {
		c, err := ConvertPoseStamped(m)
		if err != nil {
			return ret, err
		}
		ret.Poses = append(ret.Poses, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertStatisticsNames(rosMsg gengo_plotjuggler_msgs.StatisticsNames) (proto_plotjuggler_msgs.StatisticsNames, error) {

	ret := proto_plotjuggler_msgs.StatisticsNames{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert StatisticsNames, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]string, 0, len(rosMsg.Names))
	for _, m := range rosMsg.Names {
		t1 = append(t1, m)
	}

	ret.Names = append(ret.Names, t1...)

	ret.NamesVersion = rosMsg.NamesVersion

	return ret, nil
}

func ConvertStatisticsValues(rosMsg gengo_plotjuggler_msgs.StatisticsValues) (proto_plotjuggler_msgs.StatisticsValues, error) {

	ret := proto_plotjuggler_msgs.StatisticsValues{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert StatisticsValues, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float64, 0, len(rosMsg.Values))
	for _, m := range rosMsg.Values {
		t1 = append(t1, m)
	}

	ret.Values = append(ret.Values, t1...)

	ret.NamesVersion = rosMsg.NamesVersion

	return ret, nil
}

func ConvertDataPoint(rosMsg gengo_plotjuggler_msgs.DataPoint) (proto_plotjuggler_msgs.DataPoint, error) {

	ret := proto_plotjuggler_msgs.DataPoint{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert DataPoint, msg is nil")
	//}

	ret.NameIndex = uint32(rosMsg.NameIndex)

	ret.Stamp = rosMsg.Stamp

	ret.Value = rosMsg.Value

	return ret, nil
}

func ConvertDataPoints(rosMsg gengo_plotjuggler_msgs.DataPoints) (proto_plotjuggler_msgs.DataPoints, error) {

	ret := proto_plotjuggler_msgs.DataPoints{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert DataPoints, msg is nil")
	//}

	var err error

	ret.DictionaryUuid = rosMsg.DictionaryUuid

	for _, m := range rosMsg.Samples {
		c, err := ConvertDataPoint(m)
		if err != nil {
			return ret, err
		}
		ret.Samples = append(ret.Samples, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertDictionary(rosMsg gengo_plotjuggler_msgs.Dictionary) (proto_plotjuggler_msgs.Dictionary, error) {

	ret := proto_plotjuggler_msgs.Dictionary{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Dictionary, msg is nil")
	//}

	ret.DictionaryUuid = rosMsg.DictionaryUuid

	t1 := make([]string, 0, len(rosMsg.Names))
	for _, m := range rosMsg.Names {
		t1 = append(t1, m)
	}

	ret.Names = append(ret.Names, t1...)

	return ret, nil
}

func ConvertLogger(rosMsg gengo_roscpp.Logger) (proto_roscpp.Logger, error) {

	ret := proto_roscpp.Logger{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Logger, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.Level = rosMsg.Level

	return ret, nil
}

func ConvertClock(rosMsg gengo_rosgraph_msgs.Clock) (proto_rosgraph_msgs.Clock, error) {

	ret := proto_rosgraph_msgs.Clock{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Clock, msg is nil")
	//}

	return ret, nil
}

func ConvertLog(rosMsg gengo_rosgraph_msgs.Log) (proto_rosgraph_msgs.Log, error) {

	ret := proto_rosgraph_msgs.Log{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Log, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Level = uint32(rosMsg.Level)

	ret.Name = rosMsg.Name

	ret.Msg = rosMsg.Msg

	ret.File = rosMsg.File

	ret.Function = rosMsg.Function

	ret.Line = rosMsg.Line

	t7 := make([]string, 0, len(rosMsg.Topics))
	for _, m := range rosMsg.Topics {
		t7 = append(t7, m)
	}

	ret.Topics = append(ret.Topics, t7...)

	return ret, nil
}

func ConvertTopicStatistics(rosMsg gengo_rosgraph_msgs.TopicStatistics) (proto_rosgraph_msgs.TopicStatistics, error) {

	ret := proto_rosgraph_msgs.TopicStatistics{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TopicStatistics, msg is nil")
	//}

	ret.Topic = rosMsg.Topic

	ret.NodePub = rosMsg.NodePub

	ret.NodeSub = rosMsg.NodeSub

	ret.DeliveredMsgs = rosMsg.DeliveredMsgs

	ret.DroppedMsgs = rosMsg.DroppedMsgs

	ret.Traffic = rosMsg.Traffic

	return ret, nil
}

func ConvertFloats(rosMsg gengo_rospy_tutorials.Floats) (proto_rospy_tutorials.Floats, error) {

	ret := proto_rospy_tutorials.Floats{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Floats, msg is nil")
	//}

	t0 := make([]float32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		t0 = append(t0, m)
	}

	ret.Data = append(ret.Data, t0...)

	return ret, nil
}

func ConvertHeaderString(rosMsg gengo_rospy_tutorials.HeaderString) (proto_rospy_tutorials.HeaderString, error) {

	ret := proto_rospy_tutorials.HeaderString{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert HeaderString, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertTemperature(rosMsg gengo_sensor_msgs.Temperature) (proto_sensor_msgs.Temperature, error) {

	ret := proto_sensor_msgs.Temperature{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Temperature, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Temperature = rosMsg.Temperature

	ret.Variance = rosMsg.Variance

	return ret, nil
}

func ConvertBatteryState(rosMsg gengo_sensor_msgs.BatteryState) (proto_sensor_msgs.BatteryState, error) {

	ret := proto_sensor_msgs.BatteryState{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert BatteryState, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Voltage = rosMsg.Voltage

	ret.Current = rosMsg.Current

	ret.Charge = rosMsg.Charge

	ret.Capacity = rosMsg.Capacity

	ret.DesignCapacity = rosMsg.DesignCapacity

	ret.Percentage = rosMsg.Percentage

	ret.PowerSupplyStatus = uint32(rosMsg.PowerSupplyStatus)

	ret.PowerSupplyHealth = uint32(rosMsg.PowerSupplyHealth)

	ret.PowerSupplyTechnology = uint32(rosMsg.PowerSupplyTechnology)

	ret.Present = rosMsg.Present

	t11 := make([]float32, 0, len(rosMsg.CellVoltage))
	for _, m := range rosMsg.CellVoltage {
		t11 = append(t11, m)
	}

	ret.CellVoltage = append(ret.CellVoltage, t11...)

	ret.Location = rosMsg.Location

	ret.SerialNumber = rosMsg.SerialNumber

	return ret, nil
}

func ConvertFluidPressure(rosMsg gengo_sensor_msgs.FluidPressure) (proto_sensor_msgs.FluidPressure, error) {

	ret := proto_sensor_msgs.FluidPressure{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert FluidPressure, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.FluidPressure = rosMsg.FluidPressure

	ret.Variance = rosMsg.Variance

	return ret, nil
}

func ConvertMagneticField(rosMsg gengo_sensor_msgs.MagneticField) (proto_sensor_msgs.MagneticField, error) {

	ret := proto_sensor_msgs.MagneticField{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert MagneticField, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertVector3(rosMsg.MagneticField)
	ret.MagneticField = &t1

	if err != nil {
		return ret, err
	}

	t2 := make([]float64, 0, len(rosMsg.MagneticFieldCovariance))
	for _, m := range rosMsg.MagneticFieldCovariance {
		t2 = append(t2, m)
	}

	ret.MagneticFieldCovariance = append(ret.MagneticFieldCovariance, t2...)

	return ret, nil
}

func ConvertPointField(rosMsg gengo_sensor_msgs.PointField) (proto_sensor_msgs.PointField, error) {

	ret := proto_sensor_msgs.PointField{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PointField, msg is nil")
	//}

	ret.Name = rosMsg.Name

	ret.Offset = rosMsg.Offset

	ret.Datatype = uint32(rosMsg.Datatype)

	ret.Count = rosMsg.Count

	return ret, nil
}

func ConvertMultiEchoLaserScan(rosMsg gengo_sensor_msgs.MultiEchoLaserScan) (proto_sensor_msgs.MultiEchoLaserScan, error) {

	ret := proto_sensor_msgs.MultiEchoLaserScan{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert MultiEchoLaserScan, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.AngleMin = rosMsg.AngleMin

	ret.AngleMax = rosMsg.AngleMax

	ret.AngleIncrement = rosMsg.AngleIncrement

	ret.TimeIncrement = rosMsg.TimeIncrement

	ret.ScanTime = rosMsg.ScanTime

	ret.RangeMin = rosMsg.RangeMin

	ret.RangeMax = rosMsg.RangeMax

	for _, m := range rosMsg.Ranges {
		c, err := ConvertLaserEcho(m)
		if err != nil {
			return ret, err
		}
		ret.Ranges = append(ret.Ranges, &c)
	}

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Intensities {
		c, err := ConvertLaserEcho(m)
		if err != nil {
			return ret, err
		}
		ret.Intensities = append(ret.Intensities, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertPointCloud(rosMsg gengo_sensor_msgs.PointCloud) (proto_sensor_msgs.PointCloud, error) {

	ret := proto_sensor_msgs.PointCloud{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PointCloud, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Points {
		c, err := ConvertPoint32(m)
		if err != nil {
			return ret, err
		}
		ret.Points = append(ret.Points, &c)
	}

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Channels {
		c, err := ConvertChannelFloat32(m)
		if err != nil {
			return ret, err
		}
		ret.Channels = append(ret.Channels, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertRegionOfInterest(rosMsg gengo_sensor_msgs.RegionOfInterest) (proto_sensor_msgs.RegionOfInterest, error) {

	ret := proto_sensor_msgs.RegionOfInterest{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert RegionOfInterest, msg is nil")
	//}

	ret.XOffset = rosMsg.XOffset

	ret.YOffset = rosMsg.YOffset

	ret.Height = rosMsg.Height

	ret.Width = rosMsg.Width

	ret.DoRectify = rosMsg.DoRectify

	return ret, nil
}

func ConvertJointState(rosMsg gengo_sensor_msgs.JointState) (proto_sensor_msgs.JointState, error) {

	ret := proto_sensor_msgs.JointState{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert JointState, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]string, 0, len(rosMsg.Name))
	for _, m := range rosMsg.Name {
		t1 = append(t1, m)
	}

	ret.Name = append(ret.Name, t1...)

	t2 := make([]float64, 0, len(rosMsg.Position))
	for _, m := range rosMsg.Position {
		t2 = append(t2, m)
	}

	ret.Position = append(ret.Position, t2...)

	t3 := make([]float64, 0, len(rosMsg.Velocity))
	for _, m := range rosMsg.Velocity {
		t3 = append(t3, m)
	}

	ret.Velocity = append(ret.Velocity, t3...)

	t4 := make([]float64, 0, len(rosMsg.Effort))
	for _, m := range rosMsg.Effort {
		t4 = append(t4, m)
	}

	ret.Effort = append(ret.Effort, t4...)

	return ret, nil
}

func ConvertJoyFeedback(rosMsg gengo_sensor_msgs.JoyFeedback) (proto_sensor_msgs.JoyFeedback, error) {

	ret := proto_sensor_msgs.JoyFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert JoyFeedback, msg is nil")
	//}

	ret.Type = uint32(rosMsg.Type)

	ret.Id = uint32(rosMsg.Id)

	ret.Intensity = rosMsg.Intensity

	return ret, nil
}

func ConvertLaserEcho(rosMsg gengo_sensor_msgs.LaserEcho) (proto_sensor_msgs.LaserEcho, error) {

	ret := proto_sensor_msgs.LaserEcho{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert LaserEcho, msg is nil")
	//}

	t0 := make([]float32, 0, len(rosMsg.Echoes))
	for _, m := range rosMsg.Echoes {
		t0 = append(t0, m)
	}

	ret.Echoes = append(ret.Echoes, t0...)

	return ret, nil
}

func ConvertChannelFloat32(rosMsg gengo_sensor_msgs.ChannelFloat32) (proto_sensor_msgs.ChannelFloat32, error) {

	ret := proto_sensor_msgs.ChannelFloat32{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ChannelFloat32, msg is nil")
	//}

	ret.Name = rosMsg.Name

	t1 := make([]float32, 0, len(rosMsg.Values))
	for _, m := range rosMsg.Values {
		t1 = append(t1, m)
	}

	ret.Values = append(ret.Values, t1...)

	return ret, nil
}

func ConvertCompressedImage(rosMsg gengo_sensor_msgs.CompressedImage) (proto_sensor_msgs.CompressedImage, error) {

	ret := proto_sensor_msgs.CompressedImage{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert CompressedImage, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Format = rosMsg.Format

	t2 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := uint32(m)
		t2 = append(t2, a)
	}

	ret.Data = append(ret.Data, t2...)

	return ret, nil
}

func ConvertImage(rosMsg gengo_sensor_msgs.Image) (proto_sensor_msgs.Image, error) {

	ret := proto_sensor_msgs.Image{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Image, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Height = rosMsg.Height

	ret.Width = rosMsg.Width

	ret.Encoding = rosMsg.Encoding

	ret.IsBigendian = uint32(rosMsg.IsBigendian)

	ret.Step = rosMsg.Step

	t6 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := uint32(m)
		t6 = append(t6, a)
	}

	ret.Data = append(ret.Data, t6...)

	return ret, nil
}

func ConvertJoyFeedbackArray(rosMsg gengo_sensor_msgs.JoyFeedbackArray) (proto_sensor_msgs.JoyFeedbackArray, error) {

	ret := proto_sensor_msgs.JoyFeedbackArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert JoyFeedbackArray, msg is nil")
	//}

	var err error

	for _, m := range rosMsg.Array {
		c, err := ConvertJoyFeedback(m)
		if err != nil {
			return ret, err
		}
		ret.Array = append(ret.Array, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertNavSatFix(rosMsg gengo_sensor_msgs.NavSatFix) (proto_sensor_msgs.NavSatFix, error) {

	ret := proto_sensor_msgs.NavSatFix{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert NavSatFix, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertNavSatStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	ret.Latitude = rosMsg.Latitude

	ret.Longitude = rosMsg.Longitude

	ret.Altitude = rosMsg.Altitude

	t5 := make([]float64, 0, len(rosMsg.PositionCovariance))
	for _, m := range rosMsg.PositionCovariance {
		t5 = append(t5, m)
	}

	ret.PositionCovariance = append(ret.PositionCovariance, t5...)

	ret.PositionCovarianceType = uint32(rosMsg.PositionCovarianceType)

	return ret, nil
}

func ConvertRange(rosMsg gengo_sensor_msgs.Range) (proto_sensor_msgs.Range, error) {

	ret := proto_sensor_msgs.Range{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Range, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.RadiationType = uint32(rosMsg.RadiationType)

	ret.FieldOfView = rosMsg.FieldOfView

	ret.MinRange = rosMsg.MinRange

	ret.MaxRange = rosMsg.MaxRange

	ret.Range = rosMsg.Range

	return ret, nil
}

func ConvertRelativeHumidity(rosMsg gengo_sensor_msgs.RelativeHumidity) (proto_sensor_msgs.RelativeHumidity, error) {

	ret := proto_sensor_msgs.RelativeHumidity{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert RelativeHumidity, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.RelativeHumidity = rosMsg.RelativeHumidity

	ret.Variance = rosMsg.Variance

	return ret, nil
}

func ConvertPointCloud2(rosMsg gengo_sensor_msgs.PointCloud2) (proto_sensor_msgs.PointCloud2, error) {

	ret := proto_sensor_msgs.PointCloud2{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert PointCloud2, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Height = rosMsg.Height

	ret.Width = rosMsg.Width

	for _, m := range rosMsg.Fields {
		c, err := ConvertPointField(m)
		if err != nil {
			return ret, err
		}
		ret.Fields = append(ret.Fields, &c)
	}

	if err != nil {
		return ret, err
	}

	ret.IsBigendian = rosMsg.IsBigendian

	ret.PointStep = rosMsg.PointStep

	ret.RowStep = rosMsg.RowStep

	t7 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := uint32(m)
		t7 = append(t7, a)
	}

	ret.Data = append(ret.Data, t7...)

	ret.IsDense = rosMsg.IsDense

	return ret, nil
}

func ConvertTimeReference(rosMsg gengo_sensor_msgs.TimeReference) (proto_sensor_msgs.TimeReference, error) {

	ret := proto_sensor_msgs.TimeReference{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert TimeReference, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Source = rosMsg.Source

	return ret, nil
}

func ConvertImu(rosMsg gengo_sensor_msgs.Imu) (proto_sensor_msgs.Imu, error) {

	ret := proto_sensor_msgs.Imu{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Imu, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertQuaternion(rosMsg.Orientation)
	ret.Orientation = &t1

	if err != nil {
		return ret, err
	}

	t2 := make([]float64, 0, len(rosMsg.OrientationCovariance))
	for _, m := range rosMsg.OrientationCovariance {
		t2 = append(t2, m)
	}

	ret.OrientationCovariance = append(ret.OrientationCovariance, t2...)

	t3, err := ConvertVector3(rosMsg.AngularVelocity)
	ret.AngularVelocity = &t3

	if err != nil {
		return ret, err
	}

	t4 := make([]float64, 0, len(rosMsg.AngularVelocityCovariance))
	for _, m := range rosMsg.AngularVelocityCovariance {
		t4 = append(t4, m)
	}

	ret.AngularVelocityCovariance = append(ret.AngularVelocityCovariance, t4...)

	t5, err := ConvertVector3(rosMsg.LinearAcceleration)
	ret.LinearAcceleration = &t5

	if err != nil {
		return ret, err
	}

	t6 := make([]float64, 0, len(rosMsg.LinearAccelerationCovariance))
	for _, m := range rosMsg.LinearAccelerationCovariance {
		t6 = append(t6, m)
	}

	ret.LinearAccelerationCovariance = append(ret.LinearAccelerationCovariance, t6...)

	return ret, nil
}

func ConvertJoy(rosMsg gengo_sensor_msgs.Joy) (proto_sensor_msgs.Joy, error) {

	ret := proto_sensor_msgs.Joy{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Joy, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float32, 0, len(rosMsg.Axes))
	for _, m := range rosMsg.Axes {
		t1 = append(t1, m)
	}

	ret.Axes = append(ret.Axes, t1...)

	t2 := make([]int32, 0, len(rosMsg.Buttons))
	for _, m := range rosMsg.Buttons {
		t2 = append(t2, m)
	}

	ret.Buttons = append(ret.Buttons, t2...)

	return ret, nil
}

func ConvertNavSatStatus(rosMsg gengo_sensor_msgs.NavSatStatus) (proto_sensor_msgs.NavSatStatus, error) {

	ret := proto_sensor_msgs.NavSatStatus{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert NavSatStatus, msg is nil")
	//}

	ret.Status = int32(rosMsg.Status)

	ret.Service = uint32(rosMsg.Service)

	return ret, nil
}

func ConvertLaserScan(rosMsg gengo_sensor_msgs.LaserScan) (proto_sensor_msgs.LaserScan, error) {

	ret := proto_sensor_msgs.LaserScan{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert LaserScan, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.AngleMin = rosMsg.AngleMin

	ret.AngleMax = rosMsg.AngleMax

	ret.AngleIncrement = rosMsg.AngleIncrement

	ret.TimeIncrement = rosMsg.TimeIncrement

	ret.ScanTime = rosMsg.ScanTime

	ret.RangeMin = rosMsg.RangeMin

	ret.RangeMax = rosMsg.RangeMax

	t8 := make([]float32, 0, len(rosMsg.Ranges))
	for _, m := range rosMsg.Ranges {
		t8 = append(t8, m)
	}

	ret.Ranges = append(ret.Ranges, t8...)

	t9 := make([]float32, 0, len(rosMsg.Intensities))
	for _, m := range rosMsg.Intensities {
		t9 = append(t9, m)
	}

	ret.Intensities = append(ret.Intensities, t9...)

	return ret, nil
}

func ConvertMultiDOFJointState(rosMsg gengo_sensor_msgs.MultiDOFJointState) (proto_sensor_msgs.MultiDOFJointState, error) {

	ret := proto_sensor_msgs.MultiDOFJointState{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert MultiDOFJointState, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]string, 0, len(rosMsg.JointNames))
	for _, m := range rosMsg.JointNames {
		t1 = append(t1, m)
	}

	ret.JointNames = append(ret.JointNames, t1...)

	for _, m := range rosMsg.Transforms {
		c, err := ConvertTransform(m)
		if err != nil {
			return ret, err
		}
		ret.Transforms = append(ret.Transforms, &c)
	}

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Twist {
		c, err := ConvertTwist(m)
		if err != nil {
			return ret, err
		}
		ret.Twist = append(ret.Twist, &c)
	}

	if err != nil {
		return ret, err
	}

	for _, m := range rosMsg.Wrench {
		c, err := ConvertWrench(m)
		if err != nil {
			return ret, err
		}
		ret.Wrench = append(ret.Wrench, &c)
	}

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertCameraInfo(rosMsg gengo_sensor_msgs.CameraInfo) (proto_sensor_msgs.CameraInfo, error) {

	ret := proto_sensor_msgs.CameraInfo{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert CameraInfo, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Height = rosMsg.Height

	ret.Width = rosMsg.Width

	ret.DistortionModel = rosMsg.DistortionModel

	t4 := make([]float64, 0, len(rosMsg.D))
	for _, m := range rosMsg.D {
		t4 = append(t4, m)
	}

	ret.D = append(ret.D, t4...)

	t5 := make([]float64, 0, len(rosMsg.K))
	for _, m := range rosMsg.K {
		t5 = append(t5, m)
	}

	ret.K = append(ret.K, t5...)

	t6 := make([]float64, 0, len(rosMsg.R))
	for _, m := range rosMsg.R {
		t6 = append(t6, m)
	}

	ret.R = append(ret.R, t6...)

	t7 := make([]float64, 0, len(rosMsg.P))
	for _, m := range rosMsg.P {
		t7 = append(t7, m)
	}

	ret.P = append(ret.P, t7...)

	ret.BinningX = rosMsg.BinningX

	ret.BinningY = rosMsg.BinningY

	t10, err := ConvertRegionOfInterest(rosMsg.Roi)
	ret.Roi = &t10

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertIlluminance(rosMsg gengo_sensor_msgs.Illuminance) (proto_sensor_msgs.Illuminance, error) {

	ret := proto_sensor_msgs.Illuminance{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Illuminance, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Illuminance = rosMsg.Illuminance

	ret.Variance = rosMsg.Variance

	return ret, nil
}

func ConvertSmachContainerInitialStatusCmd(rosMsg gengo_smach_msgs.SmachContainerInitialStatusCmd) (proto_smach_msgs.SmachContainerInitialStatusCmd, error) {

	ret := proto_smach_msgs.SmachContainerInitialStatusCmd{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SmachContainerInitialStatusCmd, msg is nil")
	//}

	ret.Path = rosMsg.Path

	t1 := make([]string, 0, len(rosMsg.InitialStates))
	for _, m := range rosMsg.InitialStates {
		t1 = append(t1, m)
	}

	ret.InitialStates = append(ret.InitialStates, t1...)

	ret.LocalData = rosMsg.LocalData

	return ret, nil
}

func ConvertSmachContainerStatus(rosMsg gengo_smach_msgs.SmachContainerStatus) (proto_smach_msgs.SmachContainerStatus, error) {

	ret := proto_smach_msgs.SmachContainerStatus{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SmachContainerStatus, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Path = rosMsg.Path

	t2 := make([]string, 0, len(rosMsg.InitialStates))
	for _, m := range rosMsg.InitialStates {
		t2 = append(t2, m)
	}

	ret.InitialStates = append(ret.InitialStates, t2...)

	t3 := make([]string, 0, len(rosMsg.ActiveStates))
	for _, m := range rosMsg.ActiveStates {
		t3 = append(t3, m)
	}

	ret.ActiveStates = append(ret.ActiveStates, t3...)

	ret.LocalData = rosMsg.LocalData

	ret.Info = rosMsg.Info

	return ret, nil
}

func ConvertSmachContainerStructure(rosMsg gengo_smach_msgs.SmachContainerStructure) (proto_smach_msgs.SmachContainerStructure, error) {

	ret := proto_smach_msgs.SmachContainerStructure{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SmachContainerStructure, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	ret.Path = rosMsg.Path

	t2 := make([]string, 0, len(rosMsg.Children))
	for _, m := range rosMsg.Children {
		t2 = append(t2, m)
	}

	ret.Children = append(ret.Children, t2...)

	t3 := make([]string, 0, len(rosMsg.InternalOutcomes))
	for _, m := range rosMsg.InternalOutcomes {
		t3 = append(t3, m)
	}

	ret.InternalOutcomes = append(ret.InternalOutcomes, t3...)

	t4 := make([]string, 0, len(rosMsg.OutcomesFrom))
	for _, m := range rosMsg.OutcomesFrom {
		t4 = append(t4, m)
	}

	ret.OutcomesFrom = append(ret.OutcomesFrom, t4...)

	t5 := make([]string, 0, len(rosMsg.OutcomesTo))
	for _, m := range rosMsg.OutcomesTo {
		t5 = append(t5, m)
	}

	ret.OutcomesTo = append(ret.OutcomesTo, t5...)

	t6 := make([]string, 0, len(rosMsg.ContainerOutcomes))
	for _, m := range rosMsg.ContainerOutcomes {
		t6 = append(t6, m)
	}

	ret.ContainerOutcomes = append(ret.ContainerOutcomes, t6...)

	return ret, nil
}

func ConvertSoundRequestActionFeedback(rosMsg gengo_sound_play.SoundRequestActionFeedback) (proto_sound_play.SoundRequestActionFeedback, error) {

	ret := proto_sound_play.SoundRequestActionFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SoundRequestActionFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertSoundRequestFeedback(rosMsg.Feedback)
	ret.Feedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertSoundRequestActionGoal(rosMsg gengo_sound_play.SoundRequestActionGoal) (proto_sound_play.SoundRequestActionGoal, error) {

	ret := proto_sound_play.SoundRequestActionGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SoundRequestActionGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertSoundRequestGoal(rosMsg.Goal)
	ret.Goal = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertSoundRequestActionResult(rosMsg gengo_sound_play.SoundRequestActionResult) (proto_sound_play.SoundRequestActionResult, error) {

	ret := proto_sound_play.SoundRequestActionResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SoundRequestActionResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertSoundRequestResult(rosMsg.Result)
	ret.Result = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertSoundRequestFeedback(rosMsg gengo_sound_play.SoundRequestFeedback) (proto_sound_play.SoundRequestFeedback, error) {

	ret := proto_sound_play.SoundRequestFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SoundRequestFeedback, msg is nil")
	//}

	ret.Playing = rosMsg.Playing

	return ret, nil
}

func ConvertSoundRequestGoal(rosMsg gengo_sound_play.SoundRequestGoal) (proto_sound_play.SoundRequestGoal, error) {

	ret := proto_sound_play.SoundRequestGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SoundRequestGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertSoundRequest(rosMsg.SoundRequest)
	ret.SoundRequest = &t0

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertSoundRequestResult(rosMsg gengo_sound_play.SoundRequestResult) (proto_sound_play.SoundRequestResult, error) {

	ret := proto_sound_play.SoundRequestResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SoundRequestResult, msg is nil")
	//}

	ret.Playing = rosMsg.Playing

	return ret, nil
}

func ConvertSoundRequest(rosMsg gengo_sound_play.SoundRequest) (proto_sound_play.SoundRequest, error) {

	ret := proto_sound_play.SoundRequest{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SoundRequest, msg is nil")
	//}

	ret.Sound = int32(rosMsg.Sound)

	ret.Command = int32(rosMsg.Command)

	ret.Volume = rosMsg.Volume

	ret.Arg = rosMsg.Arg

	ret.Arg2 = rosMsg.Arg2

	return ret, nil
}

func ConvertSoundRequestAction(rosMsg gengo_sound_play.SoundRequestAction) (proto_sound_play.SoundRequestAction, error) {

	ret := proto_sound_play.SoundRequestAction{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert SoundRequestAction, msg is nil")
	//}

	var err error

	t0, err := ConvertSoundRequestActionGoal(rosMsg.ActionGoal)
	ret.ActionGoal = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertSoundRequestActionResult(rosMsg.ActionResult)
	ret.ActionResult = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertSoundRequestActionFeedback(rosMsg.ActionFeedback)
	ret.ActionFeedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertChar(rosMsg gengo_std_msgs.Char) (proto_std_msgs.Char, error) {

	ret := proto_std_msgs.Char{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Char, msg is nil")
	//}

	ret.Data = uint32(rosMsg.Data)

	return ret, nil
}

func ConvertInt8(rosMsg gengo_std_msgs.Int8) (proto_std_msgs.Int8, error) {

	ret := proto_std_msgs.Int8{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Int8, msg is nil")
	//}

	ret.Data = int32(rosMsg.Data)

	return ret, nil
}

func ConvertUInt16(rosMsg gengo_std_msgs.UInt16) (proto_std_msgs.UInt16, error) {

	ret := proto_std_msgs.UInt16{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UInt16, msg is nil")
	//}

	ret.Data = uint32(rosMsg.Data)

	return ret, nil
}

func ConvertFloat32MultiArray(rosMsg gengo_std_msgs.Float32MultiArray) (proto_std_msgs.Float32MultiArray, error) {

	ret := proto_std_msgs.Float32MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Float32MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		t1 = append(t1, m)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertByteMultiArray(rosMsg gengo_std_msgs.ByteMultiArray) (proto_std_msgs.ByteMultiArray, error) {

	ret := proto_std_msgs.ByteMultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ByteMultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := uint32(m)
		t1 = append(t1, a)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertFloat64MultiArray(rosMsg gengo_std_msgs.Float64MultiArray) (proto_std_msgs.Float64MultiArray, error) {

	ret := proto_std_msgs.Float64MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Float64MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]float64, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		t1 = append(t1, m)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertMultiArrayLayout(rosMsg gengo_std_msgs.MultiArrayLayout) (proto_std_msgs.MultiArrayLayout, error) {

	ret := proto_std_msgs.MultiArrayLayout{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert MultiArrayLayout, msg is nil")
	//}

	var err error

	for _, m := range rosMsg.Dim {
		c, err := ConvertMultiArrayDimension(m)
		if err != nil {
			return ret, err
		}
		ret.Dim = append(ret.Dim, &c)
	}

	if err != nil {
		return ret, err
	}

	ret.DataOffset = rosMsg.DataOffset

	return ret, nil
}

func ConvertTime(rosMsg gengo_std_msgs.Time) (proto_std_msgs.Time, error) {

	ret := proto_std_msgs.Time{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Time, msg is nil")
	//}

	return ret, nil
}

func ConvertUInt64(rosMsg gengo_std_msgs.UInt64) (proto_std_msgs.UInt64, error) {

	ret := proto_std_msgs.UInt64{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UInt64, msg is nil")
	//}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertUInt8MultiArray(rosMsg gengo_std_msgs.UInt8MultiArray) (proto_std_msgs.UInt8MultiArray, error) {

	ret := proto_std_msgs.UInt8MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UInt8MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := uint32(m)
		t1 = append(t1, a)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertFloat32(rosMsg gengo_std_msgs.Float32) (proto_std_msgs.Float32, error) {

	ret := proto_std_msgs.Float32{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Float32, msg is nil")
	//}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertInt32(rosMsg gengo_std_msgs.Int32) (proto_std_msgs.Int32, error) {

	ret := proto_std_msgs.Int32{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Int32, msg is nil")
	//}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertInt8MultiArray(rosMsg gengo_std_msgs.Int8MultiArray) (proto_std_msgs.Int8MultiArray, error) {

	ret := proto_std_msgs.Int8MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Int8MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]int32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := int32(m)
		t1 = append(t1, a)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertBool(rosMsg gengo_std_msgs.Bool) (proto_std_msgs.Bool, error) {

	ret := proto_std_msgs.Bool{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Bool, msg is nil")
	//}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertMultiArrayDimension(rosMsg gengo_std_msgs.MultiArrayDimension) (proto_std_msgs.MultiArrayDimension, error) {

	ret := proto_std_msgs.MultiArrayDimension{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert MultiArrayDimension, msg is nil")
	//}

	ret.Label = rosMsg.Label

	ret.Size = rosMsg.Size

	ret.Stride = rosMsg.Stride

	return ret, nil
}

func ConvertUInt64MultiArray(rosMsg gengo_std_msgs.UInt64MultiArray) (proto_std_msgs.UInt64MultiArray, error) {

	ret := proto_std_msgs.UInt64MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UInt64MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]uint64, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		t1 = append(t1, m)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertUInt8(rosMsg gengo_std_msgs.UInt8) (proto_std_msgs.UInt8, error) {

	ret := proto_std_msgs.UInt8{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UInt8, msg is nil")
	//}

	ret.Data = uint32(rosMsg.Data)

	return ret, nil
}

func ConvertDuration(rosMsg gengo_std_msgs.Duration) (proto_std_msgs.Duration, error) {

	ret := proto_std_msgs.Duration{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Duration, msg is nil")
	//}

	return ret, nil
}

func ConvertFloat64(rosMsg gengo_std_msgs.Float64) (proto_std_msgs.Float64, error) {

	ret := proto_std_msgs.Float64{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Float64, msg is nil")
	//}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertString(rosMsg gengo_std_msgs.String) (proto_std_msgs.String, error) {

	ret := proto_std_msgs.String{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert String, msg is nil")
	//}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertUInt16MultiArray(rosMsg gengo_std_msgs.UInt16MultiArray) (proto_std_msgs.UInt16MultiArray, error) {

	ret := proto_std_msgs.UInt16MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UInt16MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := uint32(m)
		t1 = append(t1, a)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertInt32MultiArray(rosMsg gengo_std_msgs.Int32MultiArray) (proto_std_msgs.Int32MultiArray, error) {

	ret := proto_std_msgs.Int32MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Int32MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]int32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		t1 = append(t1, m)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertInt64(rosMsg gengo_std_msgs.Int64) (proto_std_msgs.Int64, error) {

	ret := proto_std_msgs.Int64{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Int64, msg is nil")
	//}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertByte(rosMsg gengo_std_msgs.Byte) (proto_std_msgs.Byte, error) {

	ret := proto_std_msgs.Byte{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Byte, msg is nil")
	//}

	ret.Data = uint32(rosMsg.Data)

	return ret, nil
}

func ConvertColorRGBA(rosMsg gengo_std_msgs.ColorRGBA) (proto_std_msgs.ColorRGBA, error) {

	ret := proto_std_msgs.ColorRGBA{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ColorRGBA, msg is nil")
	//}

	ret.R = rosMsg.R

	ret.G = rosMsg.G

	ret.B = rosMsg.B

	ret.A = rosMsg.A

	return ret, nil
}

func ConvertEmpty(rosMsg gengo_std_msgs.Empty) (proto_std_msgs.Empty, error) {

	ret := proto_std_msgs.Empty{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Empty, msg is nil")
	//}

	return ret, nil
}

func ConvertHeader(rosMsg gengo_std_msgs.Header) (proto_std_msgs.Header, error) {

	ret := proto_std_msgs.Header{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Header, msg is nil")
	//}

	ret.Seq = rosMsg.Seq

	ret.FrameId = rosMsg.FrameId

	return ret, nil
}

func ConvertInt16(rosMsg gengo_std_msgs.Int16) (proto_std_msgs.Int16, error) {

	ret := proto_std_msgs.Int16{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Int16, msg is nil")
	//}

	ret.Data = int32(rosMsg.Data)

	return ret, nil
}

func ConvertInt16MultiArray(rosMsg gengo_std_msgs.Int16MultiArray) (proto_std_msgs.Int16MultiArray, error) {

	ret := proto_std_msgs.Int16MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Int16MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]int32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := int32(m)
		t1 = append(t1, a)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertInt64MultiArray(rosMsg gengo_std_msgs.Int64MultiArray) (proto_std_msgs.Int64MultiArray, error) {

	ret := proto_std_msgs.Int64MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Int64MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]int64, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		t1 = append(t1, m)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertUInt32(rosMsg gengo_std_msgs.UInt32) (proto_std_msgs.UInt32, error) {

	ret := proto_std_msgs.UInt32{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UInt32, msg is nil")
	//}

	ret.Data = rosMsg.Data

	return ret, nil
}

func ConvertUInt32MultiArray(rosMsg gengo_std_msgs.UInt32MultiArray) (proto_std_msgs.UInt32MultiArray, error) {

	ret := proto_std_msgs.UInt32MultiArray{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UInt32MultiArray, msg is nil")
	//}

	var err error

	t0, err := ConvertMultiArrayLayout(rosMsg.Layout)
	ret.Layout = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		t1 = append(t1, m)
	}

	ret.Data = append(ret.Data, t1...)

	return ret, nil
}

func ConvertPacket(rosMsg gengo_theora_image_transport.Packet) (proto_theora_image_transport.Packet, error) {

	ret := proto_theora_image_transport.Packet{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Packet, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1 := make([]uint32, 0, len(rosMsg.Data))
	for _, m := range rosMsg.Data {
		a := uint32(m)
		t1 = append(t1, a)
	}

	ret.Data = append(ret.Data, t1...)

	ret.BOS = rosMsg.BOS

	ret.EOS = rosMsg.EOS

	ret.Granulepos = rosMsg.Granulepos

	ret.Packetno = rosMsg.Packetno

	return ret, nil
}

func ConvertVelocity(rosMsg gengo_turtle_actionlib.Velocity) (proto_turtle_actionlib.Velocity, error) {

	ret := proto_turtle_actionlib.Velocity{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert Velocity, msg is nil")
	//}

	ret.Linear = rosMsg.Linear

	ret.Angular = rosMsg.Angular

	return ret, nil
}

func ConvertShapeAction(rosMsg gengo_turtle_actionlib.ShapeAction) (proto_turtle_actionlib.ShapeAction, error) {

	ret := proto_turtle_actionlib.ShapeAction{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ShapeAction, msg is nil")
	//}

	var err error

	t0, err := ConvertShapeActionGoal(rosMsg.ActionGoal)
	ret.ActionGoal = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertShapeActionResult(rosMsg.ActionResult)
	ret.ActionResult = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertShapeActionFeedback(rosMsg.ActionFeedback)
	ret.ActionFeedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertShapeActionFeedback(rosMsg gengo_turtle_actionlib.ShapeActionFeedback) (proto_turtle_actionlib.ShapeActionFeedback, error) {

	ret := proto_turtle_actionlib.ShapeActionFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ShapeActionFeedback, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertShapeFeedback(rosMsg.Feedback)
	ret.Feedback = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertShapeActionGoal(rosMsg gengo_turtle_actionlib.ShapeActionGoal) (proto_turtle_actionlib.ShapeActionGoal, error) {

	ret := proto_turtle_actionlib.ShapeActionGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ShapeActionGoal, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalID(rosMsg.GoalId)
	ret.GoalId = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertShapeGoal(rosMsg.Goal)
	ret.Goal = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertShapeActionResult(rosMsg gengo_turtle_actionlib.ShapeActionResult) (proto_turtle_actionlib.ShapeActionResult, error) {

	ret := proto_turtle_actionlib.ShapeActionResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ShapeActionResult, msg is nil")
	//}

	var err error

	t0, err := ConvertHeader(rosMsg.Header)
	ret.Header = &t0

	if err != nil {
		return ret, err
	}

	t1, err := ConvertGoalStatus(rosMsg.Status)
	ret.Status = &t1

	if err != nil {
		return ret, err
	}

	t2, err := ConvertShapeResult(rosMsg.Result)
	ret.Result = &t2

	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ConvertShapeFeedback(rosMsg gengo_turtle_actionlib.ShapeFeedback) (proto_turtle_actionlib.ShapeFeedback, error) {

	ret := proto_turtle_actionlib.ShapeFeedback{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ShapeFeedback, msg is nil")
	//}

	return ret, nil
}

func ConvertShapeGoal(rosMsg gengo_turtle_actionlib.ShapeGoal) (proto_turtle_actionlib.ShapeGoal, error) {

	ret := proto_turtle_actionlib.ShapeGoal{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ShapeGoal, msg is nil")
	//}

	ret.Edges = rosMsg.Edges

	ret.Radius = rosMsg.Radius

	return ret, nil
}

func ConvertShapeResult(rosMsg gengo_turtle_actionlib.ShapeResult) (proto_turtle_actionlib.ShapeResult, error) {

	ret := proto_turtle_actionlib.ShapeResult{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert ShapeResult, msg is nil")
	//}

	ret.InteriorAngle = rosMsg.InteriorAngle

	ret.Apothem = rosMsg.Apothem

	return ret, nil
}

func ConvertUniqueID(rosMsg gengo_uuid_msgs.UniqueID) (proto_uuid_msgs.UniqueID, error) {

	ret := proto_uuid_msgs.UniqueID{}

	//if rosMsg == nil {
	//  return ret, fmt.Errorf("Cannot not convert UniqueID, msg is nil")
	//}

	t0 := make([]uint32, 0, len(rosMsg.Uuid))
	for _, m := range rosMsg.Uuid {
		a := uint32(m)
		t0 = append(t0, a)
	}

	ret.Uuid = append(ret.Uuid, t0...)

	return ret, nil
}

func SerializeStateFeedback(msg proto_mpc_local_planner_msgs.StateFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeOptimalControlResult(msg proto_mpc_local_planner_msgs.OptimalControlResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePosition2DInt(msg proto_base_local_planner.Position2DInt) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeExtrinsics(msg proto_realsense2_camera.Extrinsics) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeIMUInfo(msg proto_realsense2_camera.IMUInfo) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeMetadata(msg proto_realsense2_camera.Metadata) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestFeedback(msg proto_actionlib.TestFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestRequestAction(msg proto_actionlib.TestRequestAction) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestRequestActionFeedback(msg proto_actionlib.TestRequestActionFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestRequestActionGoal(msg proto_actionlib.TestRequestActionGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwoIntsActionResult(msg proto_actionlib.TwoIntsActionResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestRequestGoal(msg proto_actionlib.TestRequestGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestAction(msg proto_actionlib.TestAction) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestActionFeedback(msg proto_actionlib.TestActionFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestRequestResult(msg proto_actionlib.TestRequestResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestResult(msg proto_actionlib.TestResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwoIntsActionGoal(msg proto_actionlib.TwoIntsActionGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwoIntsGoal(msg proto_actionlib.TwoIntsGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwoIntsResult(msg proto_actionlib.TwoIntsResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwoIntsFeedback(msg proto_actionlib.TwoIntsFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestActionGoal(msg proto_actionlib.TestActionGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestActionResult(msg proto_actionlib.TestActionResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestGoal(msg proto_actionlib.TestGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestRequestActionResult(msg proto_actionlib.TestRequestActionResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTestRequestFeedback(msg proto_actionlib.TestRequestFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwoIntsAction(msg proto_actionlib.TwoIntsAction) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwoIntsActionFeedback(msg proto_actionlib.TwoIntsActionFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGoalID(msg proto_actionlib_msgs.GoalID) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGoalStatus(msg proto_actionlib_msgs.GoalStatus) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGoalStatusArray(msg proto_actionlib_msgs.GoalStatusArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAveragingActionGoal(msg proto_actionlib_tutorials.AveragingActionGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFibonacciActionResult(msg proto_actionlib_tutorials.FibonacciActionResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFibonacciGoal(msg proto_actionlib_tutorials.FibonacciGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAveragingAction(msg proto_actionlib_tutorials.AveragingAction) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAveragingActionFeedback(msg proto_actionlib_tutorials.AveragingActionFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAveragingResult(msg proto_actionlib_tutorials.AveragingResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFibonacciFeedback(msg proto_actionlib_tutorials.FibonacciFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFibonacciResult(msg proto_actionlib_tutorials.FibonacciResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAveragingActionResult(msg proto_actionlib_tutorials.AveragingActionResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAveragingFeedback(msg proto_actionlib_tutorials.AveragingFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAveragingGoal(msg proto_actionlib_tutorials.AveragingGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFibonacciActionFeedback(msg proto_actionlib_tutorials.FibonacciActionFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFibonacciActionGoal(msg proto_actionlib_tutorials.FibonacciActionGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFibonacciAction(msg proto_actionlib_tutorials.FibonacciAction) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAudioData(msg proto_audio_common_msgs.AudioData) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAudioDataStamped(msg proto_audio_common_msgs.AudioDataStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAudioInfo(msg proto_audio_common_msgs.AudioInfo) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeConstants(msg proto_bond.Constants) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeStatus(msg proto_bond.Status) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeControllersStatistics(msg proto_controller_manager_msgs.ControllersStatistics) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeHardwareInterfaceResources(msg proto_controller_manager_msgs.HardwareInterfaceResources) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeControllerState(msg proto_controller_manager_msgs.ControllerState) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeControllerStatistics(msg proto_controller_manager_msgs.ControllerStatistics) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeKeyValue(msg proto_diagnostic_msgs.KeyValue) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeDiagnosticArray(msg proto_diagnostic_msgs.DiagnosticArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeDiagnosticStatus(msg proto_diagnostic_msgs.DiagnosticStatus) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGroupState(msg proto_dynamic_reconfigure.GroupState) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeIntParameter(msg proto_dynamic_reconfigure.IntParameter) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSensorLevels(msg proto_dynamic_reconfigure.SensorLevels) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeBoolParameter(msg proto_dynamic_reconfigure.BoolParameter) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeConfig(msg proto_dynamic_reconfigure.Config) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeConfigDescription(msg proto_dynamic_reconfigure.ConfigDescription) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeDoubleParameter(msg proto_dynamic_reconfigure.DoubleParameter) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGroup(msg proto_dynamic_reconfigure.Group) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeParamDescription(msg proto_dynamic_reconfigure.ParamDescription) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeStrParameter(msg proto_dynamic_reconfigure.StrParameter) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAccelWithCovariance(msg proto_geometry_msgs.AccelWithCovariance) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePolygon(msg proto_geometry_msgs.Polygon) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePose(msg proto_geometry_msgs.Pose) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePoseWithCovariance(msg proto_geometry_msgs.PoseWithCovariance) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTransform(msg proto_geometry_msgs.Transform) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwistStamped(msg proto_geometry_msgs.TwistStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwist(msg proto_geometry_msgs.Twist) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwistWithCovariance(msg proto_geometry_msgs.TwistWithCovariance) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePoint(msg proto_geometry_msgs.Point) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePose2D(msg proto_geometry_msgs.Pose2D) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePoseWithCovarianceStamped(msg proto_geometry_msgs.PoseWithCovarianceStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTransformStamped(msg proto_geometry_msgs.TransformStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeVector3Stamped(msg proto_geometry_msgs.Vector3Stamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeWrench(msg proto_geometry_msgs.Wrench) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAccelWithCovarianceStamped(msg proto_geometry_msgs.AccelWithCovarianceStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePoint32(msg proto_geometry_msgs.Point32) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePolygonStamped(msg proto_geometry_msgs.PolygonStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeQuaternionStamped(msg proto_geometry_msgs.QuaternionStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeQuaternion(msg proto_geometry_msgs.Quaternion) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAccelStamped(msg proto_geometry_msgs.AccelStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePointStamped(msg proto_geometry_msgs.PointStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePoseArray(msg proto_geometry_msgs.PoseArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePoseStamped(msg proto_geometry_msgs.PoseStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInertia(msg proto_geometry_msgs.Inertia) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeWrenchStamped(msg proto_geometry_msgs.WrenchStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeAccel(msg proto_geometry_msgs.Accel) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInertiaStamped(msg proto_geometry_msgs.InertiaStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTwistWithCovarianceStamped(msg proto_geometry_msgs.TwistWithCovarianceStamped) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeVector3(msg proto_geometry_msgs.Vector3) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGetMapActionFeedback(msg proto_nav_msgs.GetMapActionFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGetMapActionGoal(msg proto_nav_msgs.GetMapActionGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeMapMetaData(msg proto_nav_msgs.MapMetaData) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeOdometry(msg proto_nav_msgs.Odometry) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGetMapAction(msg proto_nav_msgs.GetMapAction) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGetMapActionResult(msg proto_nav_msgs.GetMapActionResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGetMapFeedback(msg proto_nav_msgs.GetMapFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGetMapGoal(msg proto_nav_msgs.GetMapGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGetMapResult(msg proto_nav_msgs.GetMapResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeGridCells(msg proto_nav_msgs.GridCells) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeOccupancyGrid(msg proto_nav_msgs.OccupancyGrid) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePath(msg proto_nav_msgs.Path) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeStatisticsNames(msg proto_plotjuggler_msgs.StatisticsNames) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeStatisticsValues(msg proto_plotjuggler_msgs.StatisticsValues) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeDataPoint(msg proto_plotjuggler_msgs.DataPoint) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeDataPoints(msg proto_plotjuggler_msgs.DataPoints) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeDictionary(msg proto_plotjuggler_msgs.Dictionary) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeLogger(msg proto_roscpp.Logger) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeClock(msg proto_rosgraph_msgs.Clock) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeLog(msg proto_rosgraph_msgs.Log) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTopicStatistics(msg proto_rosgraph_msgs.TopicStatistics) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFloats(msg proto_rospy_tutorials.Floats) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeHeaderString(msg proto_rospy_tutorials.HeaderString) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTemperature(msg proto_sensor_msgs.Temperature) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeBatteryState(msg proto_sensor_msgs.BatteryState) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFluidPressure(msg proto_sensor_msgs.FluidPressure) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeMagneticField(msg proto_sensor_msgs.MagneticField) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePointField(msg proto_sensor_msgs.PointField) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeMultiEchoLaserScan(msg proto_sensor_msgs.MultiEchoLaserScan) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePointCloud(msg proto_sensor_msgs.PointCloud) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeRegionOfInterest(msg proto_sensor_msgs.RegionOfInterest) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeJointState(msg proto_sensor_msgs.JointState) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeJoyFeedback(msg proto_sensor_msgs.JoyFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeLaserEcho(msg proto_sensor_msgs.LaserEcho) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeChannelFloat32(msg proto_sensor_msgs.ChannelFloat32) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeCompressedImage(msg proto_sensor_msgs.CompressedImage) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeImage(msg proto_sensor_msgs.Image) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeJoyFeedbackArray(msg proto_sensor_msgs.JoyFeedbackArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeNavSatFix(msg proto_sensor_msgs.NavSatFix) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeRange(msg proto_sensor_msgs.Range) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeRelativeHumidity(msg proto_sensor_msgs.RelativeHumidity) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePointCloud2(msg proto_sensor_msgs.PointCloud2) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTimeReference(msg proto_sensor_msgs.TimeReference) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeImu(msg proto_sensor_msgs.Imu) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeJoy(msg proto_sensor_msgs.Joy) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeNavSatStatus(msg proto_sensor_msgs.NavSatStatus) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeLaserScan(msg proto_sensor_msgs.LaserScan) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeMultiDOFJointState(msg proto_sensor_msgs.MultiDOFJointState) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeCameraInfo(msg proto_sensor_msgs.CameraInfo) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeIlluminance(msg proto_sensor_msgs.Illuminance) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSmachContainerInitialStatusCmd(msg proto_smach_msgs.SmachContainerInitialStatusCmd) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSmachContainerStatus(msg proto_smach_msgs.SmachContainerStatus) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSmachContainerStructure(msg proto_smach_msgs.SmachContainerStructure) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSoundRequestActionFeedback(msg proto_sound_play.SoundRequestActionFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSoundRequestActionGoal(msg proto_sound_play.SoundRequestActionGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSoundRequestActionResult(msg proto_sound_play.SoundRequestActionResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSoundRequestFeedback(msg proto_sound_play.SoundRequestFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSoundRequestGoal(msg proto_sound_play.SoundRequestGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSoundRequestResult(msg proto_sound_play.SoundRequestResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSoundRequest(msg proto_sound_play.SoundRequest) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeSoundRequestAction(msg proto_sound_play.SoundRequestAction) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeChar(msg proto_std_msgs.Char) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInt8(msg proto_std_msgs.Int8) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUInt16(msg proto_std_msgs.UInt16) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFloat32MultiArray(msg proto_std_msgs.Float32MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeByteMultiArray(msg proto_std_msgs.ByteMultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFloat64MultiArray(msg proto_std_msgs.Float64MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeMultiArrayLayout(msg proto_std_msgs.MultiArrayLayout) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeTime(msg proto_std_msgs.Time) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUInt64(msg proto_std_msgs.UInt64) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUInt8MultiArray(msg proto_std_msgs.UInt8MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFloat32(msg proto_std_msgs.Float32) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInt32(msg proto_std_msgs.Int32) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInt8MultiArray(msg proto_std_msgs.Int8MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeBool(msg proto_std_msgs.Bool) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeMultiArrayDimension(msg proto_std_msgs.MultiArrayDimension) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUInt64MultiArray(msg proto_std_msgs.UInt64MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUInt8(msg proto_std_msgs.UInt8) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeDuration(msg proto_std_msgs.Duration) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeFloat64(msg proto_std_msgs.Float64) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeString(msg proto_std_msgs.String) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUInt16MultiArray(msg proto_std_msgs.UInt16MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInt32MultiArray(msg proto_std_msgs.Int32MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInt64(msg proto_std_msgs.Int64) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeByte(msg proto_std_msgs.Byte) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeColorRGBA(msg proto_std_msgs.ColorRGBA) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeEmpty(msg proto_std_msgs.Empty) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeHeader(msg proto_std_msgs.Header) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInt16(msg proto_std_msgs.Int16) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInt16MultiArray(msg proto_std_msgs.Int16MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeInt64MultiArray(msg proto_std_msgs.Int64MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUInt32(msg proto_std_msgs.UInt32) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUInt32MultiArray(msg proto_std_msgs.UInt32MultiArray) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializePacket(msg proto_theora_image_transport.Packet) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeVelocity(msg proto_turtle_actionlib.Velocity) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeShapeAction(msg proto_turtle_actionlib.ShapeAction) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeShapeActionFeedback(msg proto_turtle_actionlib.ShapeActionFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeShapeActionGoal(msg proto_turtle_actionlib.ShapeActionGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeShapeActionResult(msg proto_turtle_actionlib.ShapeActionResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeShapeFeedback(msg proto_turtle_actionlib.ShapeFeedback) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeShapeGoal(msg proto_turtle_actionlib.ShapeGoal) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeShapeResult(msg proto_turtle_actionlib.ShapeResult) ([]byte, error) {
	return proto.Marshal(&msg)
}

func SerializeUniqueID(msg proto_uuid_msgs.UniqueID) ([]byte, error) {
	return proto.Marshal(&msg)
}

func GetBuilderFromName(name string) (iface.Builder, error) {
	switch name {
	case "StateFeedback":
		return &channel.BuilderUtil[gengo_mpc_local_planner_msgs.StateFeedback, proto_mpc_local_planner_msgs.StateFeedback]{
			MsgConverter: ConvertStateFeedback,
			Serializer:   SerializeStateFeedback,
		}, nil
	case "OptimalControlResult":
		return &channel.BuilderUtil[gengo_mpc_local_planner_msgs.OptimalControlResult, proto_mpc_local_planner_msgs.OptimalControlResult]{
			MsgConverter: ConvertOptimalControlResult,
			Serializer:   SerializeOptimalControlResult,
		}, nil
	case "Position2DInt":
		return &channel.BuilderUtil[gengo_base_local_planner.Position2DInt, proto_base_local_planner.Position2DInt]{
			MsgConverter: ConvertPosition2DInt,
			Serializer:   SerializePosition2DInt,
		}, nil
	case "Extrinsics":
		return &channel.BuilderUtil[gengo_realsense2_camera.Extrinsics, proto_realsense2_camera.Extrinsics]{
			MsgConverter: ConvertExtrinsics,
			Serializer:   SerializeExtrinsics,
		}, nil
	case "IMUInfo":
		return &channel.BuilderUtil[gengo_realsense2_camera.IMUInfo, proto_realsense2_camera.IMUInfo]{
			MsgConverter: ConvertIMUInfo,
			Serializer:   SerializeIMUInfo,
		}, nil
	case "Metadata":
		return &channel.BuilderUtil[gengo_realsense2_camera.Metadata, proto_realsense2_camera.Metadata]{
			MsgConverter: ConvertMetadata,
			Serializer:   SerializeMetadata,
		}, nil
	case "TestFeedback":
		return &channel.BuilderUtil[gengo_actionlib.TestFeedback, proto_actionlib.TestFeedback]{
			MsgConverter: ConvertTestFeedback,
			Serializer:   SerializeTestFeedback,
		}, nil
	case "TestRequestAction":
		return &channel.BuilderUtil[gengo_actionlib.TestRequestAction, proto_actionlib.TestRequestAction]{
			MsgConverter: ConvertTestRequestAction,
			Serializer:   SerializeTestRequestAction,
		}, nil
	case "TestRequestActionFeedback":
		return &channel.BuilderUtil[gengo_actionlib.TestRequestActionFeedback, proto_actionlib.TestRequestActionFeedback]{
			MsgConverter: ConvertTestRequestActionFeedback,
			Serializer:   SerializeTestRequestActionFeedback,
		}, nil
	case "TestRequestActionGoal":
		return &channel.BuilderUtil[gengo_actionlib.TestRequestActionGoal, proto_actionlib.TestRequestActionGoal]{
			MsgConverter: ConvertTestRequestActionGoal,
			Serializer:   SerializeTestRequestActionGoal,
		}, nil
	case "TwoIntsActionResult":
		return &channel.BuilderUtil[gengo_actionlib.TwoIntsActionResult, proto_actionlib.TwoIntsActionResult]{
			MsgConverter: ConvertTwoIntsActionResult,
			Serializer:   SerializeTwoIntsActionResult,
		}, nil
	case "TestRequestGoal":
		return &channel.BuilderUtil[gengo_actionlib.TestRequestGoal, proto_actionlib.TestRequestGoal]{
			MsgConverter: ConvertTestRequestGoal,
			Serializer:   SerializeTestRequestGoal,
		}, nil
	case "TestAction":
		return &channel.BuilderUtil[gengo_actionlib.TestAction, proto_actionlib.TestAction]{
			MsgConverter: ConvertTestAction,
			Serializer:   SerializeTestAction,
		}, nil
	case "TestActionFeedback":
		return &channel.BuilderUtil[gengo_actionlib.TestActionFeedback, proto_actionlib.TestActionFeedback]{
			MsgConverter: ConvertTestActionFeedback,
			Serializer:   SerializeTestActionFeedback,
		}, nil
	case "TestRequestResult":
		return &channel.BuilderUtil[gengo_actionlib.TestRequestResult, proto_actionlib.TestRequestResult]{
			MsgConverter: ConvertTestRequestResult,
			Serializer:   SerializeTestRequestResult,
		}, nil
	case "TestResult":
		return &channel.BuilderUtil[gengo_actionlib.TestResult, proto_actionlib.TestResult]{
			MsgConverter: ConvertTestResult,
			Serializer:   SerializeTestResult,
		}, nil
	case "TwoIntsActionGoal":
		return &channel.BuilderUtil[gengo_actionlib.TwoIntsActionGoal, proto_actionlib.TwoIntsActionGoal]{
			MsgConverter: ConvertTwoIntsActionGoal,
			Serializer:   SerializeTwoIntsActionGoal,
		}, nil
	case "TwoIntsGoal":
		return &channel.BuilderUtil[gengo_actionlib.TwoIntsGoal, proto_actionlib.TwoIntsGoal]{
			MsgConverter: ConvertTwoIntsGoal,
			Serializer:   SerializeTwoIntsGoal,
		}, nil
	case "TwoIntsResult":
		return &channel.BuilderUtil[gengo_actionlib.TwoIntsResult, proto_actionlib.TwoIntsResult]{
			MsgConverter: ConvertTwoIntsResult,
			Serializer:   SerializeTwoIntsResult,
		}, nil
	case "TwoIntsFeedback":
		return &channel.BuilderUtil[gengo_actionlib.TwoIntsFeedback, proto_actionlib.TwoIntsFeedback]{
			MsgConverter: ConvertTwoIntsFeedback,
			Serializer:   SerializeTwoIntsFeedback,
		}, nil
	case "TestActionGoal":
		return &channel.BuilderUtil[gengo_actionlib.TestActionGoal, proto_actionlib.TestActionGoal]{
			MsgConverter: ConvertTestActionGoal,
			Serializer:   SerializeTestActionGoal,
		}, nil
	case "TestActionResult":
		return &channel.BuilderUtil[gengo_actionlib.TestActionResult, proto_actionlib.TestActionResult]{
			MsgConverter: ConvertTestActionResult,
			Serializer:   SerializeTestActionResult,
		}, nil
	case "TestGoal":
		return &channel.BuilderUtil[gengo_actionlib.TestGoal, proto_actionlib.TestGoal]{
			MsgConverter: ConvertTestGoal,
			Serializer:   SerializeTestGoal,
		}, nil
	case "TestRequestActionResult":
		return &channel.BuilderUtil[gengo_actionlib.TestRequestActionResult, proto_actionlib.TestRequestActionResult]{
			MsgConverter: ConvertTestRequestActionResult,
			Serializer:   SerializeTestRequestActionResult,
		}, nil
	case "TestRequestFeedback":
		return &channel.BuilderUtil[gengo_actionlib.TestRequestFeedback, proto_actionlib.TestRequestFeedback]{
			MsgConverter: ConvertTestRequestFeedback,
			Serializer:   SerializeTestRequestFeedback,
		}, nil
	case "TwoIntsAction":
		return &channel.BuilderUtil[gengo_actionlib.TwoIntsAction, proto_actionlib.TwoIntsAction]{
			MsgConverter: ConvertTwoIntsAction,
			Serializer:   SerializeTwoIntsAction,
		}, nil
	case "TwoIntsActionFeedback":
		return &channel.BuilderUtil[gengo_actionlib.TwoIntsActionFeedback, proto_actionlib.TwoIntsActionFeedback]{
			MsgConverter: ConvertTwoIntsActionFeedback,
			Serializer:   SerializeTwoIntsActionFeedback,
		}, nil
	case "GoalID":
		return &channel.BuilderUtil[gengo_actionlib_msgs.GoalID, proto_actionlib_msgs.GoalID]{
			MsgConverter: ConvertGoalID,
			Serializer:   SerializeGoalID,
		}, nil
	case "GoalStatus":
		return &channel.BuilderUtil[gengo_actionlib_msgs.GoalStatus, proto_actionlib_msgs.GoalStatus]{
			MsgConverter: ConvertGoalStatus,
			Serializer:   SerializeGoalStatus,
		}, nil
	case "GoalStatusArray":
		return &channel.BuilderUtil[gengo_actionlib_msgs.GoalStatusArray, proto_actionlib_msgs.GoalStatusArray]{
			MsgConverter: ConvertGoalStatusArray,
			Serializer:   SerializeGoalStatusArray,
		}, nil
	case "AveragingActionGoal":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.AveragingActionGoal, proto_actionlib_tutorials.AveragingActionGoal]{
			MsgConverter: ConvertAveragingActionGoal,
			Serializer:   SerializeAveragingActionGoal,
		}, nil
	case "FibonacciActionResult":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.FibonacciActionResult, proto_actionlib_tutorials.FibonacciActionResult]{
			MsgConverter: ConvertFibonacciActionResult,
			Serializer:   SerializeFibonacciActionResult,
		}, nil
	case "FibonacciGoal":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.FibonacciGoal, proto_actionlib_tutorials.FibonacciGoal]{
			MsgConverter: ConvertFibonacciGoal,
			Serializer:   SerializeFibonacciGoal,
		}, nil
	case "AveragingAction":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.AveragingAction, proto_actionlib_tutorials.AveragingAction]{
			MsgConverter: ConvertAveragingAction,
			Serializer:   SerializeAveragingAction,
		}, nil
	case "AveragingActionFeedback":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.AveragingActionFeedback, proto_actionlib_tutorials.AveragingActionFeedback]{
			MsgConverter: ConvertAveragingActionFeedback,
			Serializer:   SerializeAveragingActionFeedback,
		}, nil
	case "AveragingResult":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.AveragingResult, proto_actionlib_tutorials.AveragingResult]{
			MsgConverter: ConvertAveragingResult,
			Serializer:   SerializeAveragingResult,
		}, nil
	case "FibonacciFeedback":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.FibonacciFeedback, proto_actionlib_tutorials.FibonacciFeedback]{
			MsgConverter: ConvertFibonacciFeedback,
			Serializer:   SerializeFibonacciFeedback,
		}, nil
	case "FibonacciResult":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.FibonacciResult, proto_actionlib_tutorials.FibonacciResult]{
			MsgConverter: ConvertFibonacciResult,
			Serializer:   SerializeFibonacciResult,
		}, nil
	case "AveragingActionResult":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.AveragingActionResult, proto_actionlib_tutorials.AveragingActionResult]{
			MsgConverter: ConvertAveragingActionResult,
			Serializer:   SerializeAveragingActionResult,
		}, nil
	case "AveragingFeedback":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.AveragingFeedback, proto_actionlib_tutorials.AveragingFeedback]{
			MsgConverter: ConvertAveragingFeedback,
			Serializer:   SerializeAveragingFeedback,
		}, nil
	case "AveragingGoal":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.AveragingGoal, proto_actionlib_tutorials.AveragingGoal]{
			MsgConverter: ConvertAveragingGoal,
			Serializer:   SerializeAveragingGoal,
		}, nil
	case "FibonacciActionFeedback":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.FibonacciActionFeedback, proto_actionlib_tutorials.FibonacciActionFeedback]{
			MsgConverter: ConvertFibonacciActionFeedback,
			Serializer:   SerializeFibonacciActionFeedback,
		}, nil
	case "FibonacciActionGoal":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.FibonacciActionGoal, proto_actionlib_tutorials.FibonacciActionGoal]{
			MsgConverter: ConvertFibonacciActionGoal,
			Serializer:   SerializeFibonacciActionGoal,
		}, nil
	case "FibonacciAction":
		return &channel.BuilderUtil[gengo_actionlib_tutorials.FibonacciAction, proto_actionlib_tutorials.FibonacciAction]{
			MsgConverter: ConvertFibonacciAction,
			Serializer:   SerializeFibonacciAction,
		}, nil
	case "AudioData":
		return &channel.BuilderUtil[gengo_audio_common_msgs.AudioData, proto_audio_common_msgs.AudioData]{
			MsgConverter: ConvertAudioData,
			Serializer:   SerializeAudioData,
		}, nil
	case "AudioDataStamped":
		return &channel.BuilderUtil[gengo_audio_common_msgs.AudioDataStamped, proto_audio_common_msgs.AudioDataStamped]{
			MsgConverter: ConvertAudioDataStamped,
			Serializer:   SerializeAudioDataStamped,
		}, nil
	case "AudioInfo":
		return &channel.BuilderUtil[gengo_audio_common_msgs.AudioInfo, proto_audio_common_msgs.AudioInfo]{
			MsgConverter: ConvertAudioInfo,
			Serializer:   SerializeAudioInfo,
		}, nil
	case "Constants":
		return &channel.BuilderUtil[gengo_bond.Constants, proto_bond.Constants]{
			MsgConverter: ConvertConstants,
			Serializer:   SerializeConstants,
		}, nil
	case "Status":
		return &channel.BuilderUtil[gengo_bond.Status, proto_bond.Status]{
			MsgConverter: ConvertStatus,
			Serializer:   SerializeStatus,
		}, nil
	case "ControllersStatistics":
		return &channel.BuilderUtil[gengo_controller_manager_msgs.ControllersStatistics, proto_controller_manager_msgs.ControllersStatistics]{
			MsgConverter: ConvertControllersStatistics,
			Serializer:   SerializeControllersStatistics,
		}, nil
	case "HardwareInterfaceResources":
		return &channel.BuilderUtil[gengo_controller_manager_msgs.HardwareInterfaceResources, proto_controller_manager_msgs.HardwareInterfaceResources]{
			MsgConverter: ConvertHardwareInterfaceResources,
			Serializer:   SerializeHardwareInterfaceResources,
		}, nil
	case "ControllerState":
		return &channel.BuilderUtil[gengo_controller_manager_msgs.ControllerState, proto_controller_manager_msgs.ControllerState]{
			MsgConverter: ConvertControllerState,
			Serializer:   SerializeControllerState,
		}, nil
	case "ControllerStatistics":
		return &channel.BuilderUtil[gengo_controller_manager_msgs.ControllerStatistics, proto_controller_manager_msgs.ControllerStatistics]{
			MsgConverter: ConvertControllerStatistics,
			Serializer:   SerializeControllerStatistics,
		}, nil
	case "KeyValue":
		return &channel.BuilderUtil[gengo_diagnostic_msgs.KeyValue, proto_diagnostic_msgs.KeyValue]{
			MsgConverter: ConvertKeyValue,
			Serializer:   SerializeKeyValue,
		}, nil
	case "DiagnosticArray":
		return &channel.BuilderUtil[gengo_diagnostic_msgs.DiagnosticArray, proto_diagnostic_msgs.DiagnosticArray]{
			MsgConverter: ConvertDiagnosticArray,
			Serializer:   SerializeDiagnosticArray,
		}, nil
	case "DiagnosticStatus":
		return &channel.BuilderUtil[gengo_diagnostic_msgs.DiagnosticStatus, proto_diagnostic_msgs.DiagnosticStatus]{
			MsgConverter: ConvertDiagnosticStatus,
			Serializer:   SerializeDiagnosticStatus,
		}, nil
	case "GroupState":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.GroupState, proto_dynamic_reconfigure.GroupState]{
			MsgConverter: ConvertGroupState,
			Serializer:   SerializeGroupState,
		}, nil
	case "IntParameter":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.IntParameter, proto_dynamic_reconfigure.IntParameter]{
			MsgConverter: ConvertIntParameter,
			Serializer:   SerializeIntParameter,
		}, nil
	case "SensorLevels":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.SensorLevels, proto_dynamic_reconfigure.SensorLevels]{
			MsgConverter: ConvertSensorLevels,
			Serializer:   SerializeSensorLevels,
		}, nil
	case "BoolParameter":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.BoolParameter, proto_dynamic_reconfigure.BoolParameter]{
			MsgConverter: ConvertBoolParameter,
			Serializer:   SerializeBoolParameter,
		}, nil
	case "Config":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.Config, proto_dynamic_reconfigure.Config]{
			MsgConverter: ConvertConfig,
			Serializer:   SerializeConfig,
		}, nil
	case "ConfigDescription":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.ConfigDescription, proto_dynamic_reconfigure.ConfigDescription]{
			MsgConverter: ConvertConfigDescription,
			Serializer:   SerializeConfigDescription,
		}, nil
	case "DoubleParameter":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.DoubleParameter, proto_dynamic_reconfigure.DoubleParameter]{
			MsgConverter: ConvertDoubleParameter,
			Serializer:   SerializeDoubleParameter,
		}, nil
	case "Group":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.Group, proto_dynamic_reconfigure.Group]{
			MsgConverter: ConvertGroup,
			Serializer:   SerializeGroup,
		}, nil
	case "ParamDescription":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.ParamDescription, proto_dynamic_reconfigure.ParamDescription]{
			MsgConverter: ConvertParamDescription,
			Serializer:   SerializeParamDescription,
		}, nil
	case "StrParameter":
		return &channel.BuilderUtil[gengo_dynamic_reconfigure.StrParameter, proto_dynamic_reconfigure.StrParameter]{
			MsgConverter: ConvertStrParameter,
			Serializer:   SerializeStrParameter,
		}, nil
	case "AccelWithCovariance":
		return &channel.BuilderUtil[gengo_geometry_msgs.AccelWithCovariance, proto_geometry_msgs.AccelWithCovariance]{
			MsgConverter: ConvertAccelWithCovariance,
			Serializer:   SerializeAccelWithCovariance,
		}, nil
	case "Polygon":
		return &channel.BuilderUtil[gengo_geometry_msgs.Polygon, proto_geometry_msgs.Polygon]{
			MsgConverter: ConvertPolygon,
			Serializer:   SerializePolygon,
		}, nil
	case "Pose":
		return &channel.BuilderUtil[gengo_geometry_msgs.Pose, proto_geometry_msgs.Pose]{
			MsgConverter: ConvertPose,
			Serializer:   SerializePose,
		}, nil
	case "PoseWithCovariance":
		return &channel.BuilderUtil[gengo_geometry_msgs.PoseWithCovariance, proto_geometry_msgs.PoseWithCovariance]{
			MsgConverter: ConvertPoseWithCovariance,
			Serializer:   SerializePoseWithCovariance,
		}, nil
	case "Transform":
		return &channel.BuilderUtil[gengo_geometry_msgs.Transform, proto_geometry_msgs.Transform]{
			MsgConverter: ConvertTransform,
			Serializer:   SerializeTransform,
		}, nil
	case "TwistStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.TwistStamped, proto_geometry_msgs.TwistStamped]{
			MsgConverter: ConvertTwistStamped,
			Serializer:   SerializeTwistStamped,
		}, nil
	case "Twist":
		return &channel.BuilderUtil[gengo_geometry_msgs.Twist, proto_geometry_msgs.Twist]{
			MsgConverter: ConvertTwist,
			Serializer:   SerializeTwist,
		}, nil
	case "TwistWithCovariance":
		return &channel.BuilderUtil[gengo_geometry_msgs.TwistWithCovariance, proto_geometry_msgs.TwistWithCovariance]{
			MsgConverter: ConvertTwistWithCovariance,
			Serializer:   SerializeTwistWithCovariance,
		}, nil
	case "Point":
		return &channel.BuilderUtil[gengo_geometry_msgs.Point, proto_geometry_msgs.Point]{
			MsgConverter: ConvertPoint,
			Serializer:   SerializePoint,
		}, nil
	case "Pose2D":
		return &channel.BuilderUtil[gengo_geometry_msgs.Pose2D, proto_geometry_msgs.Pose2D]{
			MsgConverter: ConvertPose2D,
			Serializer:   SerializePose2D,
		}, nil
	case "PoseWithCovarianceStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.PoseWithCovarianceStamped, proto_geometry_msgs.PoseWithCovarianceStamped]{
			MsgConverter: ConvertPoseWithCovarianceStamped,
			Serializer:   SerializePoseWithCovarianceStamped,
		}, nil
	case "TransformStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.TransformStamped, proto_geometry_msgs.TransformStamped]{
			MsgConverter: ConvertTransformStamped,
			Serializer:   SerializeTransformStamped,
		}, nil
	case "Vector3Stamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.Vector3Stamped, proto_geometry_msgs.Vector3Stamped]{
			MsgConverter: ConvertVector3Stamped,
			Serializer:   SerializeVector3Stamped,
		}, nil
	case "Wrench":
		return &channel.BuilderUtil[gengo_geometry_msgs.Wrench, proto_geometry_msgs.Wrench]{
			MsgConverter: ConvertWrench,
			Serializer:   SerializeWrench,
		}, nil
	case "AccelWithCovarianceStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.AccelWithCovarianceStamped, proto_geometry_msgs.AccelWithCovarianceStamped]{
			MsgConverter: ConvertAccelWithCovarianceStamped,
			Serializer:   SerializeAccelWithCovarianceStamped,
		}, nil
	case "Point32":
		return &channel.BuilderUtil[gengo_geometry_msgs.Point32, proto_geometry_msgs.Point32]{
			MsgConverter: ConvertPoint32,
			Serializer:   SerializePoint32,
		}, nil
	case "PolygonStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.PolygonStamped, proto_geometry_msgs.PolygonStamped]{
			MsgConverter: ConvertPolygonStamped,
			Serializer:   SerializePolygonStamped,
		}, nil
	case "QuaternionStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.QuaternionStamped, proto_geometry_msgs.QuaternionStamped]{
			MsgConverter: ConvertQuaternionStamped,
			Serializer:   SerializeQuaternionStamped,
		}, nil
	case "Quaternion":
		return &channel.BuilderUtil[gengo_geometry_msgs.Quaternion, proto_geometry_msgs.Quaternion]{
			MsgConverter: ConvertQuaternion,
			Serializer:   SerializeQuaternion,
		}, nil
	case "AccelStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.AccelStamped, proto_geometry_msgs.AccelStamped]{
			MsgConverter: ConvertAccelStamped,
			Serializer:   SerializeAccelStamped,
		}, nil
	case "PointStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.PointStamped, proto_geometry_msgs.PointStamped]{
			MsgConverter: ConvertPointStamped,
			Serializer:   SerializePointStamped,
		}, nil
	case "PoseArray":
		return &channel.BuilderUtil[gengo_geometry_msgs.PoseArray, proto_geometry_msgs.PoseArray]{
			MsgConverter: ConvertPoseArray,
			Serializer:   SerializePoseArray,
		}, nil
	case "PoseStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.PoseStamped, proto_geometry_msgs.PoseStamped]{
			MsgConverter: ConvertPoseStamped,
			Serializer:   SerializePoseStamped,
		}, nil
	case "Inertia":
		return &channel.BuilderUtil[gengo_geometry_msgs.Inertia, proto_geometry_msgs.Inertia]{
			MsgConverter: ConvertInertia,
			Serializer:   SerializeInertia,
		}, nil
	case "WrenchStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.WrenchStamped, proto_geometry_msgs.WrenchStamped]{
			MsgConverter: ConvertWrenchStamped,
			Serializer:   SerializeWrenchStamped,
		}, nil
	case "Accel":
		return &channel.BuilderUtil[gengo_geometry_msgs.Accel, proto_geometry_msgs.Accel]{
			MsgConverter: ConvertAccel,
			Serializer:   SerializeAccel,
		}, nil
	case "InertiaStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.InertiaStamped, proto_geometry_msgs.InertiaStamped]{
			MsgConverter: ConvertInertiaStamped,
			Serializer:   SerializeInertiaStamped,
		}, nil
	case "TwistWithCovarianceStamped":
		return &channel.BuilderUtil[gengo_geometry_msgs.TwistWithCovarianceStamped, proto_geometry_msgs.TwistWithCovarianceStamped]{
			MsgConverter: ConvertTwistWithCovarianceStamped,
			Serializer:   SerializeTwistWithCovarianceStamped,
		}, nil
	case "Vector3":
		return &channel.BuilderUtil[gengo_geometry_msgs.Vector3, proto_geometry_msgs.Vector3]{
			MsgConverter: ConvertVector3,
			Serializer:   SerializeVector3,
		}, nil
	case "GetMapActionFeedback":
		return &channel.BuilderUtil[gengo_nav_msgs.GetMapActionFeedback, proto_nav_msgs.GetMapActionFeedback]{
			MsgConverter: ConvertGetMapActionFeedback,
			Serializer:   SerializeGetMapActionFeedback,
		}, nil
	case "GetMapActionGoal":
		return &channel.BuilderUtil[gengo_nav_msgs.GetMapActionGoal, proto_nav_msgs.GetMapActionGoal]{
			MsgConverter: ConvertGetMapActionGoal,
			Serializer:   SerializeGetMapActionGoal,
		}, nil
	case "MapMetaData":
		return &channel.BuilderUtil[gengo_nav_msgs.MapMetaData, proto_nav_msgs.MapMetaData]{
			MsgConverter: ConvertMapMetaData,
			Serializer:   SerializeMapMetaData,
		}, nil
	case "Odometry":
		return &channel.BuilderUtil[gengo_nav_msgs.Odometry, proto_nav_msgs.Odometry]{
			MsgConverter: ConvertOdometry,
			Serializer:   SerializeOdometry,
		}, nil
	case "GetMapAction":
		return &channel.BuilderUtil[gengo_nav_msgs.GetMapAction, proto_nav_msgs.GetMapAction]{
			MsgConverter: ConvertGetMapAction,
			Serializer:   SerializeGetMapAction,
		}, nil
	case "GetMapActionResult":
		return &channel.BuilderUtil[gengo_nav_msgs.GetMapActionResult, proto_nav_msgs.GetMapActionResult]{
			MsgConverter: ConvertGetMapActionResult,
			Serializer:   SerializeGetMapActionResult,
		}, nil
	case "GetMapFeedback":
		return &channel.BuilderUtil[gengo_nav_msgs.GetMapFeedback, proto_nav_msgs.GetMapFeedback]{
			MsgConverter: ConvertGetMapFeedback,
			Serializer:   SerializeGetMapFeedback,
		}, nil
	case "GetMapGoal":
		return &channel.BuilderUtil[gengo_nav_msgs.GetMapGoal, proto_nav_msgs.GetMapGoal]{
			MsgConverter: ConvertGetMapGoal,
			Serializer:   SerializeGetMapGoal,
		}, nil
	case "GetMapResult":
		return &channel.BuilderUtil[gengo_nav_msgs.GetMapResult, proto_nav_msgs.GetMapResult]{
			MsgConverter: ConvertGetMapResult,
			Serializer:   SerializeGetMapResult,
		}, nil
	case "GridCells":
		return &channel.BuilderUtil[gengo_nav_msgs.GridCells, proto_nav_msgs.GridCells]{
			MsgConverter: ConvertGridCells,
			Serializer:   SerializeGridCells,
		}, nil
	case "OccupancyGrid":
		return &channel.BuilderUtil[gengo_nav_msgs.OccupancyGrid, proto_nav_msgs.OccupancyGrid]{
			MsgConverter: ConvertOccupancyGrid,
			Serializer:   SerializeOccupancyGrid,
		}, nil
	case "Path":
		return &channel.BuilderUtil[gengo_nav_msgs.Path, proto_nav_msgs.Path]{
			MsgConverter: ConvertPath,
			Serializer:   SerializePath,
		}, nil
	case "StatisticsNames":
		return &channel.BuilderUtil[gengo_plotjuggler_msgs.StatisticsNames, proto_plotjuggler_msgs.StatisticsNames]{
			MsgConverter: ConvertStatisticsNames,
			Serializer:   SerializeStatisticsNames,
		}, nil
	case "StatisticsValues":
		return &channel.BuilderUtil[gengo_plotjuggler_msgs.StatisticsValues, proto_plotjuggler_msgs.StatisticsValues]{
			MsgConverter: ConvertStatisticsValues,
			Serializer:   SerializeStatisticsValues,
		}, nil
	case "DataPoint":
		return &channel.BuilderUtil[gengo_plotjuggler_msgs.DataPoint, proto_plotjuggler_msgs.DataPoint]{
			MsgConverter: ConvertDataPoint,
			Serializer:   SerializeDataPoint,
		}, nil
	case "DataPoints":
		return &channel.BuilderUtil[gengo_plotjuggler_msgs.DataPoints, proto_plotjuggler_msgs.DataPoints]{
			MsgConverter: ConvertDataPoints,
			Serializer:   SerializeDataPoints,
		}, nil
	case "Dictionary":
		return &channel.BuilderUtil[gengo_plotjuggler_msgs.Dictionary, proto_plotjuggler_msgs.Dictionary]{
			MsgConverter: ConvertDictionary,
			Serializer:   SerializeDictionary,
		}, nil
	case "Logger":
		return &channel.BuilderUtil[gengo_roscpp.Logger, proto_roscpp.Logger]{
			MsgConverter: ConvertLogger,
			Serializer:   SerializeLogger,
		}, nil
	case "Clock":
		return &channel.BuilderUtil[gengo_rosgraph_msgs.Clock, proto_rosgraph_msgs.Clock]{
			MsgConverter: ConvertClock,
			Serializer:   SerializeClock,
		}, nil
	case "Log":
		return &channel.BuilderUtil[gengo_rosgraph_msgs.Log, proto_rosgraph_msgs.Log]{
			MsgConverter: ConvertLog,
			Serializer:   SerializeLog,
		}, nil
	case "TopicStatistics":
		return &channel.BuilderUtil[gengo_rosgraph_msgs.TopicStatistics, proto_rosgraph_msgs.TopicStatistics]{
			MsgConverter: ConvertTopicStatistics,
			Serializer:   SerializeTopicStatistics,
		}, nil
	case "Floats":
		return &channel.BuilderUtil[gengo_rospy_tutorials.Floats, proto_rospy_tutorials.Floats]{
			MsgConverter: ConvertFloats,
			Serializer:   SerializeFloats,
		}, nil
	case "HeaderString":
		return &channel.BuilderUtil[gengo_rospy_tutorials.HeaderString, proto_rospy_tutorials.HeaderString]{
			MsgConverter: ConvertHeaderString,
			Serializer:   SerializeHeaderString,
		}, nil
	case "Temperature":
		return &channel.BuilderUtil[gengo_sensor_msgs.Temperature, proto_sensor_msgs.Temperature]{
			MsgConverter: ConvertTemperature,
			Serializer:   SerializeTemperature,
		}, nil
	case "BatteryState":
		return &channel.BuilderUtil[gengo_sensor_msgs.BatteryState, proto_sensor_msgs.BatteryState]{
			MsgConverter: ConvertBatteryState,
			Serializer:   SerializeBatteryState,
		}, nil
	case "FluidPressure":
		return &channel.BuilderUtil[gengo_sensor_msgs.FluidPressure, proto_sensor_msgs.FluidPressure]{
			MsgConverter: ConvertFluidPressure,
			Serializer:   SerializeFluidPressure,
		}, nil
	case "MagneticField":
		return &channel.BuilderUtil[gengo_sensor_msgs.MagneticField, proto_sensor_msgs.MagneticField]{
			MsgConverter: ConvertMagneticField,
			Serializer:   SerializeMagneticField,
		}, nil
	case "PointField":
		return &channel.BuilderUtil[gengo_sensor_msgs.PointField, proto_sensor_msgs.PointField]{
			MsgConverter: ConvertPointField,
			Serializer:   SerializePointField,
		}, nil
	case "MultiEchoLaserScan":
		return &channel.BuilderUtil[gengo_sensor_msgs.MultiEchoLaserScan, proto_sensor_msgs.MultiEchoLaserScan]{
			MsgConverter: ConvertMultiEchoLaserScan,
			Serializer:   SerializeMultiEchoLaserScan,
		}, nil
	case "PointCloud":
		return &channel.BuilderUtil[gengo_sensor_msgs.PointCloud, proto_sensor_msgs.PointCloud]{
			MsgConverter: ConvertPointCloud,
			Serializer:   SerializePointCloud,
		}, nil
	case "RegionOfInterest":
		return &channel.BuilderUtil[gengo_sensor_msgs.RegionOfInterest, proto_sensor_msgs.RegionOfInterest]{
			MsgConverter: ConvertRegionOfInterest,
			Serializer:   SerializeRegionOfInterest,
		}, nil
	case "JointState":
		return &channel.BuilderUtil[gengo_sensor_msgs.JointState, proto_sensor_msgs.JointState]{
			MsgConverter: ConvertJointState,
			Serializer:   SerializeJointState,
		}, nil
	case "JoyFeedback":
		return &channel.BuilderUtil[gengo_sensor_msgs.JoyFeedback, proto_sensor_msgs.JoyFeedback]{
			MsgConverter: ConvertJoyFeedback,
			Serializer:   SerializeJoyFeedback,
		}, nil
	case "LaserEcho":
		return &channel.BuilderUtil[gengo_sensor_msgs.LaserEcho, proto_sensor_msgs.LaserEcho]{
			MsgConverter: ConvertLaserEcho,
			Serializer:   SerializeLaserEcho,
		}, nil
	case "ChannelFloat32":
		return &channel.BuilderUtil[gengo_sensor_msgs.ChannelFloat32, proto_sensor_msgs.ChannelFloat32]{
			MsgConverter: ConvertChannelFloat32,
			Serializer:   SerializeChannelFloat32,
		}, nil
	case "CompressedImage":
		return &channel.BuilderUtil[gengo_sensor_msgs.CompressedImage, proto_sensor_msgs.CompressedImage]{
			MsgConverter: ConvertCompressedImage,
			Serializer:   SerializeCompressedImage,
		}, nil
	case "Image":
		return &channel.BuilderUtil[gengo_sensor_msgs.Image, proto_sensor_msgs.Image]{
			MsgConverter: ConvertImage,
			Serializer:   SerializeImage,
		}, nil
	case "JoyFeedbackArray":
		return &channel.BuilderUtil[gengo_sensor_msgs.JoyFeedbackArray, proto_sensor_msgs.JoyFeedbackArray]{
			MsgConverter: ConvertJoyFeedbackArray,
			Serializer:   SerializeJoyFeedbackArray,
		}, nil
	case "NavSatFix":
		return &channel.BuilderUtil[gengo_sensor_msgs.NavSatFix, proto_sensor_msgs.NavSatFix]{
			MsgConverter: ConvertNavSatFix,
			Serializer:   SerializeNavSatFix,
		}, nil
	case "Range":
		return &channel.BuilderUtil[gengo_sensor_msgs.Range, proto_sensor_msgs.Range]{
			MsgConverter: ConvertRange,
			Serializer:   SerializeRange,
		}, nil
	case "RelativeHumidity":
		return &channel.BuilderUtil[gengo_sensor_msgs.RelativeHumidity, proto_sensor_msgs.RelativeHumidity]{
			MsgConverter: ConvertRelativeHumidity,
			Serializer:   SerializeRelativeHumidity,
		}, nil
	case "PointCloud2":
		return &channel.BuilderUtil[gengo_sensor_msgs.PointCloud2, proto_sensor_msgs.PointCloud2]{
			MsgConverter: ConvertPointCloud2,
			Serializer:   SerializePointCloud2,
		}, nil
	case "TimeReference":
		return &channel.BuilderUtil[gengo_sensor_msgs.TimeReference, proto_sensor_msgs.TimeReference]{
			MsgConverter: ConvertTimeReference,
			Serializer:   SerializeTimeReference,
		}, nil
	case "Imu":
		return &channel.BuilderUtil[gengo_sensor_msgs.Imu, proto_sensor_msgs.Imu]{
			MsgConverter: ConvertImu,
			Serializer:   SerializeImu,
		}, nil
	case "Joy":
		return &channel.BuilderUtil[gengo_sensor_msgs.Joy, proto_sensor_msgs.Joy]{
			MsgConverter: ConvertJoy,
			Serializer:   SerializeJoy,
		}, nil
	case "NavSatStatus":
		return &channel.BuilderUtil[gengo_sensor_msgs.NavSatStatus, proto_sensor_msgs.NavSatStatus]{
			MsgConverter: ConvertNavSatStatus,
			Serializer:   SerializeNavSatStatus,
		}, nil
	case "LaserScan":
		return &channel.BuilderUtil[gengo_sensor_msgs.LaserScan, proto_sensor_msgs.LaserScan]{
			MsgConverter: ConvertLaserScan,
			Serializer:   SerializeLaserScan,
		}, nil
	case "MultiDOFJointState":
		return &channel.BuilderUtil[gengo_sensor_msgs.MultiDOFJointState, proto_sensor_msgs.MultiDOFJointState]{
			MsgConverter: ConvertMultiDOFJointState,
			Serializer:   SerializeMultiDOFJointState,
		}, nil
	case "CameraInfo":
		return &channel.BuilderUtil[gengo_sensor_msgs.CameraInfo, proto_sensor_msgs.CameraInfo]{
			MsgConverter: ConvertCameraInfo,
			Serializer:   SerializeCameraInfo,
		}, nil
	case "Illuminance":
		return &channel.BuilderUtil[gengo_sensor_msgs.Illuminance, proto_sensor_msgs.Illuminance]{
			MsgConverter: ConvertIlluminance,
			Serializer:   SerializeIlluminance,
		}, nil
	case "SmachContainerInitialStatusCmd":
		return &channel.BuilderUtil[gengo_smach_msgs.SmachContainerInitialStatusCmd, proto_smach_msgs.SmachContainerInitialStatusCmd]{
			MsgConverter: ConvertSmachContainerInitialStatusCmd,
			Serializer:   SerializeSmachContainerInitialStatusCmd,
		}, nil
	case "SmachContainerStatus":
		return &channel.BuilderUtil[gengo_smach_msgs.SmachContainerStatus, proto_smach_msgs.SmachContainerStatus]{
			MsgConverter: ConvertSmachContainerStatus,
			Serializer:   SerializeSmachContainerStatus,
		}, nil
	case "SmachContainerStructure":
		return &channel.BuilderUtil[gengo_smach_msgs.SmachContainerStructure, proto_smach_msgs.SmachContainerStructure]{
			MsgConverter: ConvertSmachContainerStructure,
			Serializer:   SerializeSmachContainerStructure,
		}, nil
	case "SoundRequestActionFeedback":
		return &channel.BuilderUtil[gengo_sound_play.SoundRequestActionFeedback, proto_sound_play.SoundRequestActionFeedback]{
			MsgConverter: ConvertSoundRequestActionFeedback,
			Serializer:   SerializeSoundRequestActionFeedback,
		}, nil
	case "SoundRequestActionGoal":
		return &channel.BuilderUtil[gengo_sound_play.SoundRequestActionGoal, proto_sound_play.SoundRequestActionGoal]{
			MsgConverter: ConvertSoundRequestActionGoal,
			Serializer:   SerializeSoundRequestActionGoal,
		}, nil
	case "SoundRequestActionResult":
		return &channel.BuilderUtil[gengo_sound_play.SoundRequestActionResult, proto_sound_play.SoundRequestActionResult]{
			MsgConverter: ConvertSoundRequestActionResult,
			Serializer:   SerializeSoundRequestActionResult,
		}, nil
	case "SoundRequestFeedback":
		return &channel.BuilderUtil[gengo_sound_play.SoundRequestFeedback, proto_sound_play.SoundRequestFeedback]{
			MsgConverter: ConvertSoundRequestFeedback,
			Serializer:   SerializeSoundRequestFeedback,
		}, nil
	case "SoundRequestGoal":
		return &channel.BuilderUtil[gengo_sound_play.SoundRequestGoal, proto_sound_play.SoundRequestGoal]{
			MsgConverter: ConvertSoundRequestGoal,
			Serializer:   SerializeSoundRequestGoal,
		}, nil
	case "SoundRequestResult":
		return &channel.BuilderUtil[gengo_sound_play.SoundRequestResult, proto_sound_play.SoundRequestResult]{
			MsgConverter: ConvertSoundRequestResult,
			Serializer:   SerializeSoundRequestResult,
		}, nil
	case "SoundRequest":
		return &channel.BuilderUtil[gengo_sound_play.SoundRequest, proto_sound_play.SoundRequest]{
			MsgConverter: ConvertSoundRequest,
			Serializer:   SerializeSoundRequest,
		}, nil
	case "SoundRequestAction":
		return &channel.BuilderUtil[gengo_sound_play.SoundRequestAction, proto_sound_play.SoundRequestAction]{
			MsgConverter: ConvertSoundRequestAction,
			Serializer:   SerializeSoundRequestAction,
		}, nil
	case "Char":
		return &channel.BuilderUtil[gengo_std_msgs.Char, proto_std_msgs.Char]{
			MsgConverter: ConvertChar,
			Serializer:   SerializeChar,
		}, nil
	case "Int8":
		return &channel.BuilderUtil[gengo_std_msgs.Int8, proto_std_msgs.Int8]{
			MsgConverter: ConvertInt8,
			Serializer:   SerializeInt8,
		}, nil
	case "UInt16":
		return &channel.BuilderUtil[gengo_std_msgs.UInt16, proto_std_msgs.UInt16]{
			MsgConverter: ConvertUInt16,
			Serializer:   SerializeUInt16,
		}, nil
	case "Float32MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.Float32MultiArray, proto_std_msgs.Float32MultiArray]{
			MsgConverter: ConvertFloat32MultiArray,
			Serializer:   SerializeFloat32MultiArray,
		}, nil
	case "ByteMultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.ByteMultiArray, proto_std_msgs.ByteMultiArray]{
			MsgConverter: ConvertByteMultiArray,
			Serializer:   SerializeByteMultiArray,
		}, nil
	case "Float64MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.Float64MultiArray, proto_std_msgs.Float64MultiArray]{
			MsgConverter: ConvertFloat64MultiArray,
			Serializer:   SerializeFloat64MultiArray,
		}, nil
	case "MultiArrayLayout":
		return &channel.BuilderUtil[gengo_std_msgs.MultiArrayLayout, proto_std_msgs.MultiArrayLayout]{
			MsgConverter: ConvertMultiArrayLayout,
			Serializer:   SerializeMultiArrayLayout,
		}, nil
	case "Time":
		return &channel.BuilderUtil[gengo_std_msgs.Time, proto_std_msgs.Time]{
			MsgConverter: ConvertTime,
			Serializer:   SerializeTime,
		}, nil
	case "UInt64":
		return &channel.BuilderUtil[gengo_std_msgs.UInt64, proto_std_msgs.UInt64]{
			MsgConverter: ConvertUInt64,
			Serializer:   SerializeUInt64,
		}, nil
	case "UInt8MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.UInt8MultiArray, proto_std_msgs.UInt8MultiArray]{
			MsgConverter: ConvertUInt8MultiArray,
			Serializer:   SerializeUInt8MultiArray,
		}, nil
	case "Float32":
		return &channel.BuilderUtil[gengo_std_msgs.Float32, proto_std_msgs.Float32]{
			MsgConverter: ConvertFloat32,
			Serializer:   SerializeFloat32,
		}, nil
	case "Int32":
		return &channel.BuilderUtil[gengo_std_msgs.Int32, proto_std_msgs.Int32]{
			MsgConverter: ConvertInt32,
			Serializer:   SerializeInt32,
		}, nil
	case "Int8MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.Int8MultiArray, proto_std_msgs.Int8MultiArray]{
			MsgConverter: ConvertInt8MultiArray,
			Serializer:   SerializeInt8MultiArray,
		}, nil
	case "Bool":
		return &channel.BuilderUtil[gengo_std_msgs.Bool, proto_std_msgs.Bool]{
			MsgConverter: ConvertBool,
			Serializer:   SerializeBool,
		}, nil
	case "MultiArrayDimension":
		return &channel.BuilderUtil[gengo_std_msgs.MultiArrayDimension, proto_std_msgs.MultiArrayDimension]{
			MsgConverter: ConvertMultiArrayDimension,
			Serializer:   SerializeMultiArrayDimension,
		}, nil
	case "UInt64MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.UInt64MultiArray, proto_std_msgs.UInt64MultiArray]{
			MsgConverter: ConvertUInt64MultiArray,
			Serializer:   SerializeUInt64MultiArray,
		}, nil
	case "UInt8":
		return &channel.BuilderUtil[gengo_std_msgs.UInt8, proto_std_msgs.UInt8]{
			MsgConverter: ConvertUInt8,
			Serializer:   SerializeUInt8,
		}, nil
	case "Duration":
		return &channel.BuilderUtil[gengo_std_msgs.Duration, proto_std_msgs.Duration]{
			MsgConverter: ConvertDuration,
			Serializer:   SerializeDuration,
		}, nil
	case "Float64":
		return &channel.BuilderUtil[gengo_std_msgs.Float64, proto_std_msgs.Float64]{
			MsgConverter: ConvertFloat64,
			Serializer:   SerializeFloat64,
		}, nil
	case "String":
		return &channel.BuilderUtil[gengo_std_msgs.String, proto_std_msgs.String]{
			MsgConverter: ConvertString,
			Serializer:   SerializeString,
		}, nil
	case "UInt16MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.UInt16MultiArray, proto_std_msgs.UInt16MultiArray]{
			MsgConverter: ConvertUInt16MultiArray,
			Serializer:   SerializeUInt16MultiArray,
		}, nil
	case "Int32MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.Int32MultiArray, proto_std_msgs.Int32MultiArray]{
			MsgConverter: ConvertInt32MultiArray,
			Serializer:   SerializeInt32MultiArray,
		}, nil
	case "Int64":
		return &channel.BuilderUtil[gengo_std_msgs.Int64, proto_std_msgs.Int64]{
			MsgConverter: ConvertInt64,
			Serializer:   SerializeInt64,
		}, nil
	case "Byte":
		return &channel.BuilderUtil[gengo_std_msgs.Byte, proto_std_msgs.Byte]{
			MsgConverter: ConvertByte,
			Serializer:   SerializeByte,
		}, nil
	case "ColorRGBA":
		return &channel.BuilderUtil[gengo_std_msgs.ColorRGBA, proto_std_msgs.ColorRGBA]{
			MsgConverter: ConvertColorRGBA,
			Serializer:   SerializeColorRGBA,
		}, nil
	case "Empty":
		return &channel.BuilderUtil[gengo_std_msgs.Empty, proto_std_msgs.Empty]{
			MsgConverter: ConvertEmpty,
			Serializer:   SerializeEmpty,
		}, nil
	case "Header":
		return &channel.BuilderUtil[gengo_std_msgs.Header, proto_std_msgs.Header]{
			MsgConverter: ConvertHeader,
			Serializer:   SerializeHeader,
		}, nil
	case "Int16":
		return &channel.BuilderUtil[gengo_std_msgs.Int16, proto_std_msgs.Int16]{
			MsgConverter: ConvertInt16,
			Serializer:   SerializeInt16,
		}, nil
	case "Int16MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.Int16MultiArray, proto_std_msgs.Int16MultiArray]{
			MsgConverter: ConvertInt16MultiArray,
			Serializer:   SerializeInt16MultiArray,
		}, nil
	case "Int64MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.Int64MultiArray, proto_std_msgs.Int64MultiArray]{
			MsgConverter: ConvertInt64MultiArray,
			Serializer:   SerializeInt64MultiArray,
		}, nil
	case "UInt32":
		return &channel.BuilderUtil[gengo_std_msgs.UInt32, proto_std_msgs.UInt32]{
			MsgConverter: ConvertUInt32,
			Serializer:   SerializeUInt32,
		}, nil
	case "UInt32MultiArray":
		return &channel.BuilderUtil[gengo_std_msgs.UInt32MultiArray, proto_std_msgs.UInt32MultiArray]{
			MsgConverter: ConvertUInt32MultiArray,
			Serializer:   SerializeUInt32MultiArray,
		}, nil
	case "Packet":
		return &channel.BuilderUtil[gengo_theora_image_transport.Packet, proto_theora_image_transport.Packet]{
			MsgConverter: ConvertPacket,
			Serializer:   SerializePacket,
		}, nil
	case "Velocity":
		return &channel.BuilderUtil[gengo_turtle_actionlib.Velocity, proto_turtle_actionlib.Velocity]{
			MsgConverter: ConvertVelocity,
			Serializer:   SerializeVelocity,
		}, nil
	case "ShapeAction":
		return &channel.BuilderUtil[gengo_turtle_actionlib.ShapeAction, proto_turtle_actionlib.ShapeAction]{
			MsgConverter: ConvertShapeAction,
			Serializer:   SerializeShapeAction,
		}, nil
	case "ShapeActionFeedback":
		return &channel.BuilderUtil[gengo_turtle_actionlib.ShapeActionFeedback, proto_turtle_actionlib.ShapeActionFeedback]{
			MsgConverter: ConvertShapeActionFeedback,
			Serializer:   SerializeShapeActionFeedback,
		}, nil
	case "ShapeActionGoal":
		return &channel.BuilderUtil[gengo_turtle_actionlib.ShapeActionGoal, proto_turtle_actionlib.ShapeActionGoal]{
			MsgConverter: ConvertShapeActionGoal,
			Serializer:   SerializeShapeActionGoal,
		}, nil
	case "ShapeActionResult":
		return &channel.BuilderUtil[gengo_turtle_actionlib.ShapeActionResult, proto_turtle_actionlib.ShapeActionResult]{
			MsgConverter: ConvertShapeActionResult,
			Serializer:   SerializeShapeActionResult,
		}, nil
	case "ShapeFeedback":
		return &channel.BuilderUtil[gengo_turtle_actionlib.ShapeFeedback, proto_turtle_actionlib.ShapeFeedback]{
			MsgConverter: ConvertShapeFeedback,
			Serializer:   SerializeShapeFeedback,
		}, nil
	case "ShapeGoal":
		return &channel.BuilderUtil[gengo_turtle_actionlib.ShapeGoal, proto_turtle_actionlib.ShapeGoal]{
			MsgConverter: ConvertShapeGoal,
			Serializer:   SerializeShapeGoal,
		}, nil
	case "ShapeResult":
		return &channel.BuilderUtil[gengo_turtle_actionlib.ShapeResult, proto_turtle_actionlib.ShapeResult]{
			MsgConverter: ConvertShapeResult,
			Serializer:   SerializeShapeResult,
		}, nil
	case "UniqueID":
		return &channel.BuilderUtil[gengo_uuid_msgs.UniqueID, proto_uuid_msgs.UniqueID]{
			MsgConverter: ConvertUniqueID,
			Serializer:   SerializeUniqueID,
		}, nil
	default:
		return nil, fmt.Errorf("Unrecognized Type")
	}
}
func AssignBuilder() bool {
	utils.NewBuilder("ros-rmq", GetBuilderFromName)
	return true
}

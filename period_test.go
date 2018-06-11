package calendar

import "testing"

var r7in19s = []struct {
	year int
	leap bool
}{
	{-359, true}, {-349, false}, {-339, false}, {-329, true}, {-319, false}, {-309, false},
	{-358, false}, {-348, true}, {-338, false}, {-328, false}, {-318, true}, {-308, false},
	{-357, false}, {-347, false}, {-337, true}, {-327, false}, {-317, false}, {-307, true},
	{-356, true}, {-346, false}, {-336, false}, {-326, true}, {-316, false}, {-306, false},
	{-355, false}, {-345, true}, {-335, false}, {-325, false}, {-315, true}, {-305, false},
	{-364, true}, {-354, false}, {-344, false}, {-334, true}, {-324, false}, {-314, false}, {-304, true},
	{-363, false}, {-353, true}, {-343, false}, {-333, false}, {-323, true}, {-313, false}, {-303, false},
	{-362, false}, {-352, false}, {-342, true}, {-332, false}, {-322, false}, {-312, true}, {-302, true},
	{-361, true}, {-351, false}, {-341, false}, {-331, true}, {-321, true}, {-311, false}, {-301, false},
	{-360, false}, {-350, true}, {-340, true}, {-330, false}, {-320, false}, {-310, true}, {-300, false},
	{-299, true}, {-289, false}, {-279, false}, {-269, true}, {-259, false}, {-249, false}, {-239, true}, {-229, false}, {-219, false}, {-209, true},
	{-298, false}, {-288, true}, {-278, false}, {-268, false}, {-258, true}, {-248, false}, {-238, false}, {-228, true}, {-218, false}, {-208, false},
	{-297, false}, {-287, false}, {-277, true}, {-267, false}, {-257, false}, {-247, true}, {-237, false}, {-227, false}, {-217, true}, {-207, true},
	{-296, true}, {-286, false}, {-276, false}, {-266, true}, {-256, false}, {-246, false}, {-236, true}, {-226, true}, {-216, false}, {-206, false},
	{-295, false}, {-285, true}, {-275, false}, {-265, false}, {-255, true}, {-245, true}, {-235, false}, {-225, false}, {-215, true}, {-205, false},
	{-294, false}, {-284, false}, {-274, true}, {-264, true}, {-254, false}, {-244, false}, {-234, true}, {-224, false}, {-214, false}, {-204, true},
	{-293, true}, {-283, true}, {-273, false}, {-263, false}, {-253, true}, {-243, false}, {-233, false}, {-223, true}, {-213, false}, {-203, false},
	{-292, false}, {-282, false}, {-272, true}, {-262, false}, {-252, false}, {-242, true}, {-232, false}, {-222, false}, {-212, true}, {-202, false},
	{-291, true}, {-281, false}, {-271, false}, {-261, true}, {-251, false}, {-241, false}, {-231, true}, {-221, false}, {-211, false}, {-201, true},
	{-290, false}, {-280, true}, {-270, false}, {-260, false}, {-250, true}, {-240, false}, {-230, false}, {-220, true}, {-210, false}, {-200, false},
	{-199, false}, {-189, false}, {-179, true}, {-169, true}, {-159, false}, {-149, false}, {-139, true}, {-129, false}, {-119, false}, {-109, true},
	{-198, true}, {-188, true}, {-178, false}, {-168, false}, {-158, true}, {-148, false}, {-138, false}, {-128, true}, {-118, false}, {-108, false},
	{-197, false}, {-187, false}, {-177, true}, {-167, false}, {-157, false}, {-147, true}, {-137, false}, {-127, false}, {-117, true}, {-107, false},
	{-196, true}, {-186, false}, {-176, false}, {-166, true}, {-156, false}, {-146, false}, {-136, true}, {-126, false}, {-116, false}, {-106, true},
	{-195, false}, {-185, true}, {-175, false}, {-165, false}, {-155, true}, {-145, false}, {-135, false}, {-125, true}, {-115, false}, {-105, false},
	{-194, false}, {-184, false}, {-174, true}, {-164, false}, {-154, false}, {-144, true}, {-134, false}, {-124, false}, {-114, true}, {-104, true},
	{-193, true}, {-183, false}, {-173, false}, {-163, true}, {-153, false}, {-143, false}, {-133, true}, {-123, true}, {-113, false}, {-103, false},
	{-192, false}, {-182, true}, {-172, false}, {-162, false}, {-152, true}, {-142, true}, {-132, false}, {-122, false}, {-112, true},
	{-191, false}, {-181, false}, {-171, true}, {-161, true}, {-151, false}, {-141, false}, {-131, true}, {-121, false}, {-111, false},
	{-190, true}, {-180, false}, {-170, false}, {-160, false}, {-150, true}, {-140, false}, {-130, false}, {-120, true}, {-110, false},
}

func TestChkR7in19(t *testing.T) {
	for _, tc := range r7in19s {
		if leap := chkR7in19(tc.year); leap != tc.leap {
			t.Errorf("year %d got %v, expect %v", tc.year, leap, tc.leap)
		}
	}
}

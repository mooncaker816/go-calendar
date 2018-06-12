package calendar

import "testing"

var r7in19s = []struct {
	year int
	leap bool
}{
	{-721, false}, {-720, true},
	{-719, false}, {-718, false}, {-717, true}, {-716, false}, {-715, false}, {-714, true}, {-713, false}, {-712, true}, {-711, false}, {-710, false}, {-709, true}, {-708, false}, {-707, false}, {-706, true}, {-705, false}, {-704, false}, {-703, true}, {-702, false}, {-701, true},
	{-700, false}, {-699, false}, {-698, true}, {-697, false}, {-696, false}, {-695, true}, {-694, false}, {-693, true}, {-692, false}, {-691, false}, {-690, true}, {-689, false}, {-688, false}, {-687, true}, {-686, false}, {-685, false}, {-684, true}, {-683, false}, {-682, true},
	{-681, false}, {-680, false}, {-679, true}, {-678, false}, {-677, false}, {-676, true}, {-675, false}, {-674, true}, {-673, false}, {-672, false}, {-671, true}, {-670, false}, {-669, false}, {-668, true}, {-667, false}, {-666, false}, {-665, true}, {-664, false}, {-663, true},
	{-662, false}, {-661, false}, {-660, true}, {-659, false}, {-658, false}, {-657, true}, {-656, false}, {-655, true}, {-654, false}, {-653, false}, {-652, true}, {-651, false}, {-650, false}, {-649, true}, {-648, false}, {-647, false}, {-646, true}, {-645, false}, {-644, true},
	{-643, false}, {-642, false}, {-641, true}, {-640, false}, {-639, false}, {-638, true}, {-637, false}, {-636, true}, {-635, false}, {-634, false}, {-633, true}, {-632, false}, {-631, false}, {-630, true}, {-629, false}, {-628, false}, {-627, true}, {-626, false}, {-625, true},
	{-624, false}, {-623, false}, {-622, true}, {-621, false}, {-620, false}, {-619, true}, {-618, false}, {-617, true}, {-616, false}, {-615, false}, {-614, true}, {-613, false}, {-612, false}, {-611, true}, {-610, false}, {-609, false}, {-608, true}, {-607, false}, {-606, true},
	{-605, false}, {-604, false}, {-603, true}, {-602, false}, {-601, false}, {-600, true}, {-599, false}, {-598, true}, {-597, false}, {-596, false}, {-595, true}, {-594, false}, {-593, false}, {-592, true}, {-591, false}, {-590, false}, {-589, true}, {-588, false}, {-587, true},
	{-586, false}, {-585, false}, {-584, true}, {-583, false}, {-582, false}, {-581, true}, {-580, false}, {-579, true}, {-578, false}, {-577, false}, {-576, true}, {-575, false}, {-574, false}, {-573, true}, {-572, false}, {-571, false}, {-570, true}, {-569, false}, {-568, true},
	{-567, false}, {-566, false}, {-565, true}, {-564, false}, {-563, false}, {-562, true}, {-561, false}, {-560, true}, {-559, false}, {-558, false}, {-557, true}, {-556, false}, {-555, false}, {-554, true}, {-553, false}, {-552, false}, {-551, true}, {-550, false}, {-549, true},
	{-548, false}, {-547, false}, {-546, true}, {-545, false}, {-544, false}, {-543, true}, {-542, false}, {-541, true}, {-540, false}, {-539, false}, {-538, true}, {-537, false}, {-536, false}, {-535, true}, {-534, false}, {-533, false}, {-532, true}, {-531, false}, {-530, true},
	{-529, false}, {-528, false}, {-527, true}, {-526, false}, {-525, false}, {-524, true}, {-523, false}, {-522, true}, {-521, false}, {-520, false}, {-519, true}, {-518, false}, {-517, false}, {-516, true}, {-515, false}, {-514, false}, {-513, true}, {-512, false}, {-511, true},
	{-510, false}, {-509, false}, {-508, true}, {-507, false}, {-506, false}, {-505, true}, {-504, false}, {-503, true}, {-502, false}, {-501, false}, {-500, true}, {-499, false}, {-498, false}, {-497, true}, {-496, false}, {-495, false}, {-494, true}, {-493, false}, {-492, true},
	{-491, false}, {-490, false}, {-489, true}, {-488, false}, {-487, false}, {-486, true}, {-485, false}, {-484, true}, {-483, false}, {-482, false}, {-481, true}, {-480, false}, {-479, false}, {-478, true}, {-477, false}, {-476, false}, {-475, true}, {-474, false}, {-473, true},
	{-472, false}, {-471, false}, {-470, true}, {-469, false}, {-468, false}, {-467, true}, {-466, false}, {-465, true}, {-464, false}, {-463, false}, {-462, true}, {-461, false}, {-460, false}, {-459, true}, {-458, false}, {-457, false}, {-456, true}, {-455, false}, {-454, true},
	{-453, false}, {-452, false}, {-451, true}, {-450, false}, {-449, false}, {-448, true}, {-447, false}, {-446, true}, {-445, false}, {-444, false}, {-443, true}, {-442, false}, {-441, false}, {-440, true}, {-439, false}, {-438, false}, {-437, true}, {-436, false}, {-435, true},
	{-434, false}, {-433, false}, {-432, true}, {-431, false}, {-430, false}, {-429, true}, {-428, false}, {-427, true}, {-426, false}, {-425, false}, {-424, true}, {-423, false}, {-422, false}, {-421, true}, {-420, false}, {-419, false}, {-418, true}, {-417, false}, {-416, true},
	{-415, false}, {-414, false}, {-413, true}, {-412, false}, {-411, false}, {-410, true}, {-409, false}, {-408, true}, {-407, false}, {-406, false}, {-405, true}, {-404, false}, {-403, false}, {-402, true}, {-401, false}, {-400, false}, {-399, true}, {-398, false}, {-397, true},
	{-396, false}, {-395, false}, {-394, true}, {-393, false}, {-392, false}, {-391, true}, {-390, false}, {-389, true}, {-388, false}, {-387, false}, {-386, true}, {-385, false}, {-384, false}, {-383, true}, {-382, false}, {-381, false}, {-380, true}, {-379, false}, {-378, true},
	{-377, false}, {-376, false}, {-375, true}, {-374, false}, {-373, false}, {-372, true}, {-371, false}, {-370, true}, {-369, false}, {-368, false}, {-367, true}, {-366, false}, {-365, false}, {-364, true}, {-363, false}, {-362, false}, {-361, true}, {-360, false}, {-359, true},
	{-358, false}, {-357, false}, {-356, true}, {-355, false}, {-354, false}, {-353, true}, {-352, false}, {-351, true}, {-350, false}, {-349, false}, {-348, true}, {-347, false}, {-346, false}, {-345, true}, {-344, false}, {-343, false}, {-342, true}, {-341, false}, {-340, true},
	{-339, false}, {-338, false}, {-337, true}, {-336, false}, {-335, false}, {-334, true}, {-333, false}, {-332, true}, {-331, false}, {-330, false}, {-329, true}, {-328, false}, {-327, false}, {-326, true}, {-325, false}, {-324, false}, {-323, true}, {-322, false}, {-321, true},
	{-320, false}, {-319, false}, {-318, true}, {-317, false}, {-316, false}, {-315, true}, {-314, false}, {-313, true}, {-312, false}, {-311, false}, {-310, true}, {-309, false}, {-308, false}, {-307, true}, {-306, false}, {-305, false}, {-304, true}, {-303, false}, {-302, true},
	{-301, false}, {-300, false}, {-299, true}, {-298, false}, {-297, false}, {-296, true}, {-295, false}, {-294, true}, {-293, false}, {-292, false}, {-291, true}, {-290, false}, {-289, false}, {-288, true}, {-287, false}, {-286, false}, {-285, true}, {-284, false}, {-283, true},
	{-282, false}, {-281, false}, {-280, true}, {-279, false}, {-278, false}, {-277, true}, {-276, false}, {-275, true}, {-274, false}, {-273, false}, {-272, true}, {-271, false}, {-270, false}, {-269, true}, {-268, false}, {-267, false}, {-266, true}, {-265, false}, {-264, true},
	{-263, false}, {-262, false}, {-261, true}, {-260, false}, {-259, false}, {-258, true}, {-257, false}, {-256, true}, {-255, false}, {-254, false}, {-253, true}, {-252, false}, {-251, false}, {-250, true}, {-249, false}, {-248, false}, {-247, true}, {-246, false}, {-245, true},
	{-244, false}, {-243, false}, {-242, true}, {-241, false}, {-240, false}, {-239, true}, {-238, false}, {-237, true}, {-236, false}, {-235, false}, {-234, true}, {-233, false}, {-232, false}, {-231, true}, {-230, false}, {-229, false}, {-228, true}, {-227, false}, {-226, true},
	{-225, false}, {-224, false}, {-223, true}, {-222, false}, {-221, false}, {-220, true}, {-219, false}, {-218, false}, {-217, true}, {-216, false}, {-215, true}, {-214, false}, {-213, false}, {-212, true}, {-211, false}, {-210, false}, {-209, true}, {-208, false}, {-207, true},
	{-206, false}, {-205, false}, {-204, true}, {-203, false}, {-202, false}, {-201, true}, {-200, false}, {-199, false}, {-198, true}, {-197, false}, {-196, true}, {-195, false}, {-194, false}, {-193, true}, {-192, false}, {-191, false}, {-190, true}, {-189, false}, {-188, true},
	{-187, false}, {-186, false}, {-185, true}, {-184, false}, {-183, false}, {-182, true}, {-181, false}, {-180, false}, {-179, true}, {-178, false}, {-177, true}, {-176, false}, {-175, false}, {-174, true}, {-173, false}, {-172, false}, {-171, true}, {-170, false}, {-169, true},
	{-168, false}, {-167, false}, {-166, true}, {-165, false}, {-164, false}, {-163, true}, {-162, false}, {-161, true}, {-160, false}, {-159, false}, {-158, true}, {-157, false}, {-156, false}, {-155, true}, {-154, false}, {-153, false}, {-152, true}, {-151, false}, {-150, true},
	{-149, false}, {-148, false}, {-147, true}, {-146, false}, {-145, false}, {-144, true}, {-143, false}, {-142, true}, {-141, false}, {-140, false}, {-139, true}, {-138, false}, {-137, false}, {-136, true}, {-135, false}, {-134, false}, {-133, true}, {-132, false}, {-131, true},
	{-130, false}, {-129, false}, {-128, true}, {-127, false}, {-126, false}, {-125, true}, {-124, false}, {-123, true}, {-122, false}, {-121, false}, {-120, true}, {-119, false}, {-118, false}, {-117, true}, {-116, false}, {-115, false}, {-114, true}, {-113, false}, {-112, true},
	{-111, false}, {-110, false}, {-109, true}, {-108, false}, {-107, false}, {-106, true}, {-105, false}, {-104, true},
}

func TestChkR7in19(t *testing.T) {
	for _, tc := range r7in19s {
		if leap := chkR7in19(tc.year); leap != tc.leap {
			t.Errorf("year %d got %v, expect %v", tc.year, leap, tc.leap)
		}
	}
}

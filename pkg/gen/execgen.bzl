# Generated by genbzl

EXECGEN_SRCS = [
  "//pkg/col/coldata:vec.eg.go",
  "//pkg/sql/colconv:datum_to_vec.eg.go",
  "//pkg/sql/colconv:vec_to_datum.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_any_not_null_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_avg_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_bool_and_or_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_concat_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_count_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_default_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_min_max_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_sum_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:hash_sum_int_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_any_not_null_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_avg_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_bool_and_or_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_concat_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_count_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_default_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_min_max_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_sum_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:ordered_sum_int_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:window_avg_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:window_bool_and_or_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:window_concat_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:window_count_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:window_min_max_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:window_sum_agg.eg.go",
  "//pkg/sql/colexec/colexecagg:window_sum_int_agg.eg.go",
  "//pkg/sql/colexec/colexecbase:cast.eg.go",
  "//pkg/sql/colexec/colexecbase:const.eg.go",
  "//pkg/sql/colexec/colexecbase:distinct.eg.go",
  "//pkg/sql/colexec/colexeccmp:default_cmp_expr.eg.go",
  "//pkg/sql/colexec/colexechash:hash_utils.eg.go",
  "//pkg/sql/colexec/colexechash:hashtable_distinct.eg.go",
  "//pkg/sql/colexec/colexechash:hashtable_full_default.eg.go",
  "//pkg/sql/colexec/colexechash:hashtable_full_deleting.eg.go",
  "//pkg/sql/colexec/colexecjoin:crossjoiner.eg.go",
  "//pkg/sql/colexec/colexecjoin:hashjoiner.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoinbase.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_exceptall.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_fullouter.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_inner.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_intersectall.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_leftanti.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_leftouter.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_leftsemi.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_rightanti.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_rightouter.eg.go",
  "//pkg/sql/colexec/colexecjoin:mergejoiner_rightsemi.eg.go",
  "//pkg/sql/colexec/colexecproj:default_cmp_proj_op.eg.go",
  "//pkg/sql/colexec/colexecproj:proj_non_const_ops.eg.go",
  "//pkg/sql/colexec/colexecprojconst:default_cmp_proj_const_op.eg.go",
  "//pkg/sql/colexec/colexecprojconst:proj_const_left_ops.eg.go",
  "//pkg/sql/colexec/colexecprojconst:proj_const_right_ops.eg.go",
  "//pkg/sql/colexec/colexecprojconst:proj_like_ops.eg.go",
  "//pkg/sql/colexec/colexecsel:default_cmp_sel_ops.eg.go",
  "//pkg/sql/colexec/colexecsel:sel_like_ops.eg.go",
  "//pkg/sql/colexec/colexecsel:selection_ops.eg.go",
  "//pkg/sql/colexec/colexecspan:span_encoder.eg.go",
  "//pkg/sql/colexec/colexecwindow:first_value.eg.go",
  "//pkg/sql/colexec/colexecwindow:lag.eg.go",
  "//pkg/sql/colexec/colexecwindow:last_value.eg.go",
  "//pkg/sql/colexec/colexecwindow:lead.eg.go",
  "//pkg/sql/colexec/colexecwindow:min_max_removable_agg.eg.go",
  "//pkg/sql/colexec/colexecwindow:nth_value.eg.go",
  "//pkg/sql/colexec/colexecwindow:ntile.eg.go",
  "//pkg/sql/colexec/colexecwindow:range_offset_handler.eg.go",
  "//pkg/sql/colexec/colexecwindow:rank.eg.go",
  "//pkg/sql/colexec/colexecwindow:relative_rank.eg.go",
  "//pkg/sql/colexec/colexecwindow:row_number.eg.go",
  "//pkg/sql/colexec/colexecwindow:window_aggregator.eg.go",
  "//pkg/sql/colexec/colexecwindow:window_framer.eg.go",
  "//pkg/sql/colexec/colexecwindow:window_peer_grouper.eg.go",
  "//pkg/sql/colexec:and_or_projection.eg.go",
  "//pkg/sql/colexec:hash_aggregator.eg.go",
  "//pkg/sql/colexec:is_null_ops.eg.go",
  "//pkg/sql/colexec:ordered_synchronizer.eg.go",
  "//pkg/sql/colexec:pdqsort.eg.go",
  "//pkg/sql/colexec:rowtovec.eg.go",
  "//pkg/sql/colexec:select_in.eg.go",
  "//pkg/sql/colexec:sort.eg.go",
  "//pkg/sql/colexec:sort_partitioner.eg.go",
  "//pkg/sql/colexec:sorttopk.eg.go",
  "//pkg/sql/colexec:substring.eg.go",
  "//pkg/sql/colexec:values_differ.eg.go",
  "//pkg/sql/colexec:vec_comparators.eg.go",
]

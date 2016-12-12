terraform remote config \
-backend=s3 \
-backend-config="bucket=${var.environment}" \
-backend-config="region=${var.root_var}" \
-backend-config="key=${var.test_derived}"

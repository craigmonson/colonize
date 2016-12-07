terraform remote config \
	-backend=s3 \
	-backend-config="region=${region}" \
	-backend-config="bucket=terraform-states-${var.environment}" \
	-backend-config="key=${var.tmpl_base_path}_${var.basename}_${var.environment}.tfstate"

module Dapp
  module Build
    module Stage
      # Source3
      class Source3 < SourceBase
        def initialize(application, next_stage)
          @prev_stage = InfraSetup.new(application, self)
          super
        end

        protected

        def dependencies_checksum
          hashsum [super,
                   setup_dependencies_files_checksum,
                   *application.builder.setup_checksum]
        end

        private

        def setup_dependencies_files_checksum
          @setup_files_checksum ||= dependencies_files_checksum(application.config._setup_dependencies)
        end
      end # Source3
    end # Stage
  end # Build
end # Dapp

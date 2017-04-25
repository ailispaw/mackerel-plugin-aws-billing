# A dummy plugin for Barge to set hostname and network correctly at the very first `vagrant up`
module VagrantPlugins
  module GuestLinux
    class Plugin < Vagrant.plugin("2")
      guest_capability("linux", "change_host_name") { Cap::ChangeHostName }
      guest_capability("linux", "configure_networks") { Cap::ConfigureNetworks }
    end
  end
end

Vagrant.configure(2) do |config|
  config.vm.define "mackerel-plugin-aws-billing"

  config.vm.box = "ailispaw/barge"

  config.vm.synced_folder ".", "/vagrant"

  config.vm.provision :docker do |d|
    d.run "mackerel-plugin-aws-billing",
      image: "ailispaw/mackerel-plugin-aws-billing",
      args: [
        "--env-file /vagrant/.env",
      ].join(" "),
      restart: false
  end

  config.vm.provision :shell do |sh|
    sh.privileged = false
    sh.inline = <<-EOT
      cat /vagrant/crontab | crontab -;  crontab -l
    EOT
  end
end

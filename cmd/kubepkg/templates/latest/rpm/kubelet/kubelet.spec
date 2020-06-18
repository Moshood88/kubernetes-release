Name: kubelet
Version: {{ .Version }}
Release: {{ .Revision }}
Summary: Container cluster management

License: ASL 2.0
URL: https://kubernetes.io
Source0: {{ .DownloadLinkBase }}/bin/linux/{{ .GoArch }}/kubelet
Source1: {{ .CNIDownloadLink }}

BuildRequires: systemd
BuildRequires: curl
Requires: iptables >= 1.4.21
Requires: socat
Requires: util-linux
Requires: ethtool
Requires: iproute
Requires: ebtables
Requires: conntrack
Conflicts: kubernetes-cni

%description
The node agent of Kubernetes, the container cluster manager.

%prep
cp -p %SOURCE0 %{_builddir}/
%setup -c -D -T -a 5 -n cni-plugins

# TODO: Do we need these?
#%autosetup
#%build
#%configure
#%make_build

%install
# TODO: Do we need this?
#rm -rf $RPM_BUILD_ROOT

cd %{_builddir}
install -m 755 -d %{buildroot}%{_unitdir}
install -m 755 -d %{buildroot}%{_unitdir}/kubelet.service.d/
install -m 755 -d %{buildroot}%{_bindir}
install -m 755 -d %{buildroot}/var/lib/kubelet/
install -p -m 755 -t %{buildroot}%{_bindir}/ kubelet
install -p -m 644 -t %{buildroot}%{_unitdir}/ kubelet.service
install -m 755 -d %{buildroot}%{_sysconfdir}/sysconfig/
install -p -m 644 -T kubelet.env %{buildroot}%{_sysconfdir}/sysconfig/kubelet
install -m 755 -d %{buildroot}%{_sysconfdir}/cni/net.d/
install -m 755 -d %{buildroot}/opt/cni/bin
mv cni-plugins/* %{buildroot}/opt/cni/bin/

# TODO: Do we need this?
#%make_install

%files
%{_bindir}/kubelet
%{_unitdir}/kubelet.service
%{_sysconfdir}/kubernetes/manifests/
/opt/cni

%config(noreplace) %{_sysconfdir}/sysconfig/kubelet

# TODO: Do we need these?
#%license add-license-file-here
#%doc add-docs-here


%changelog
* Thu Jun 18 2020 Sascha Grunert <sgrunert@suse.com.com> - 1.18.4
- Do not Obsolete with kubernetes-cni package to unbreak previous
  releases

* Fri May 22 2020 Stephen Augustus <saugustus@vmware.com> - 1.18.4
- Bundle CNI plugins (v0.8.6) in kubelet package

* Sat Jan  4 2020 Stephen Augustus <saugustus@vmware.com> - 1.18.0
- Move kubeadm into separate spec file
- Move kubectl into separate spec file
- Move kubernetes-cni into separate spec file
- Move cri-tools into separate spec file

package io.octohub.service;

import io.octohub.entity.User;
import io.octohub.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Collection;
import java.util.List;

@Service
public class UserDetailsServiceImpl implements UserDetailsService {
    
    @Autowired
    private UserRepository userRepository;
    
    @Override
    @Transactional
    public UserDetails loadUserByUsername(String email) throws UsernameNotFoundException {
        User user = userRepository.findByEmail(email)
                .orElseThrow(() -> new UsernameNotFoundException("用户未找到: " + email));
        
        return UserPrincipal.create(user);
    }

    @Transactional
    public UserDetails loadUserByUserId(Long id) throws UsernameNotFoundException {
        User user = userRepository.findById(id)
                .orElseThrow(() -> new UsernameNotFoundException("用户未找到: " + id));
        
        return UserPrincipal.create(user);
    }
    
    // 内部类：用户主体
    public static class UserPrincipal implements UserDetails {
        private Long id;
        private String email;
        private String password;
        
        public UserPrincipal(Long id, String email, String password) {
            this.id = id;
            this.email = email;
            this.password = password;
        }
        
        public static UserPrincipal create(User user) {        
            return new UserPrincipal(
                    user.getId(),
                    user.getEmail(),
                    user.getPassword()
            );
        }
        
        public Long getId() {
            return id;
        }
        
        @Override
        public String getUsername() {
            return email;
        }
        
        @Override
        public String getPassword() {
            return password;
        }
        
        @Override
        public boolean isAccountNonExpired() {
            return true;
        }
        
        @Override
        public boolean isAccountNonLocked() {
            return true;
        }
        
        @Override
        public boolean isCredentialsNonExpired() {
            return true;
        }
        
        @Override
        public boolean isEnabled() {
            return true;
        }
        
        @Override
        public Collection<? extends GrantedAuthority> getAuthorities() {
            return List.of(new SimpleGrantedAuthority("ROLE_USER"));
        }
    
    }
} 
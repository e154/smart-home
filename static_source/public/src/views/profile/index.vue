<template>
  <div class="app-container">
    <el-empty :image-size="200"></el-empty>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { UserModule } from '@/store/modules/user'

export interface IProfile {
  name: string
  email: string
  avatar: string
  roles: string
}

const defaultProfile: IProfile = {
  name: 'Loading...',
  email: 'Loading...',
  avatar: 'Loading...',
  roles: 'Loading...'
}

@Component({
  name: 'Profile',
  components: {

  }
})
export default class extends Vue {
  private user = defaultProfile
  private activeTab = 'activity'

  get name() {
    return UserModule.name
  }

  get email() {
    return UserModule.email
  }

  get avatar() {
    return UserModule.avatar
  }

  get roles() {
    return UserModule.roles
  }

  created() {
    this.getUser()
  }

  private getUser() {
    this.user = {
      name: this.name,
      email: this.email,
      avatar: this.avatar,
      roles: this.roles.join(' | ')
    }
  }
}
</script>

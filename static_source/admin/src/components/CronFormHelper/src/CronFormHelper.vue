<script setup lang="ts">
</script>

<template>
  <div class="tip custom-block w-[100%]">
    <div>
      <div>The cron expression is made of six fields. Each field can have the following
        values.
      </div>
      <table class="dYaBhr">
        <thead>
        <tr>
          <th>*</th>
          <th>*</th>
          <th>*</th>
          <th>*</th>
          <th>*</th>
          <th>*</th>
        </tr>
        </thead>
        <tbody>
        <tr>
          <td>second (0-59)</td>
          <td>minute (0-59)</td>
          <td>hour (0 - 23)</td>
          <td>day of the month (1 - 31)</td>
          <td>month (1 - 12)</td>
          <td>day of the week (0 - 6)</td>
        </tr>
        </tbody>
      </table>
      <div>Here are some examples for you.</div>
      <table class="dYaBhr">
        <thead>
        <tr>
          <th>Cron expression</th>
          <th>Schedule</th>
        </tr>
        </thead>
        <tbody>
        <tr>
          <td>* * * * * *</td>
          <td>Every second</td>
        </tr>
        <tr>
          <td>0 * * * * *</td>
          <td>Every minute</td>
        </tr>
        <tr>
          <td>0 0 * * * *</td>
          <td>Every hour</td>
        </tr>
        <tr>
          <td>0 0 0 * * *</td>
          <td>Every day at 12:00 AM</td>
        </tr>
        <tr>
          <td>0 0 0 * * FRI</td>
          <td>At 12:00 AM, only on Friday</td>
        </tr>
        <tr>
          <td>0 0 0 1 * *</td>
          <td>At 12:00 AM, on day 1 of the month</td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style lang="less" scoped>
.light {
  .custom-block.tip {
    padding: 8px 16px;
    background-color: #409eff1a;
    border-radius: 4px;
    border-left: 5px solid var(--el-color-primary);
  }
}

.dark {
  .custom-block.tip {
    padding: 8px 16px;
    background-color: #409eff1a;
    border-radius: 4px;
    border-left: 5px solid var(--el-color-primary);
  }
}

.dYaBhr thead th {
  border: 1px solid ;
  padding: 0 10px;
}

.dYaBhr tbody td {
  border: 1px solid ;
  padding: 0 10px;
}

</style>

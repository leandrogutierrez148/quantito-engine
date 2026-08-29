-- ============================================================================
-- QUANTITO ENGINE - MySQL Database Schema & Stored Procedures
-- Engine: InnoDB | Character Set: utf8mb4 | Collation: utf8mb4_unicode_ci
-- ============================================================================

CREATE DATABASE IF NOT EXISTS `quantito-engine`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `quantito-engine`;

-- Temporarily disable foreign key checks for clean recreation if necessary
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `asset_actions`;
DROP TABLE IF EXISTS `asset_variants`;
DROP TABLE IF EXISTS `actions`;
DROP TABLE IF EXISTS `implementations`;
DROP TABLE IF EXISTS `variants`;
DROP TABLE IF EXISTS `strategies`;
DROP TABLE IF EXISTS `experiment_assets`;
DROP TABLE IF EXISTS `experiments`;

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 1. TABLE: EXPERIMENTS (Quantitative Experiments)
-- ============================================================================
CREATE TABLE `experiments` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `experiment` VARCHAR(100) NOT NULL,
  `description` VARCHAR(255) NULL,
  `params` JSON NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- 'ACTIVE', 'PAUSED', 'COMPLETED'
  `notes` TEXT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_experiments_status` (`status`),
  INDEX `idx_experiments_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 2. TABLE: EXPERIMENT_ASSETS (Asset Assignment to Experiment)
-- ============================================================================
CREATE TABLE `experiment_assets` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `experiment_id` INT UNSIGNED NOT NULL,
  `asset_id` INT UNSIGNED NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- 'ACTIVE', 'EXCLUDED'
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_experiment_asset` (`experiment_id`, `asset_id`),
  INDEX `idx_experiment_assets_asset` (`asset_id`, `status`),
  INDEX `idx_experiment_assets_exp` (`experiment_id`, `status`),
  CONSTRAINT `fk_exp_assets_experiment`
    FOREIGN KEY (`experiment_id`) REFERENCES `experiments` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 3. TABLE: STRATEGIES (Quantitative Strategies Catalog)
-- ============================================================================
CREATE TABLE `strategies` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `class` VARCHAR(100) NOT NULL, -- e.g.: 'fifty_day_avg', 'price_momentum'
  `params` JSON NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_strategies_class` (`class`),
  INDEX `idx_strategies_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 4. TABLE: VARIANTS (Experiment Variants / Strategy Configurations)
-- ============================================================================
CREATE TABLE `variants` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `experiment_id` INT UNSIGNED NOT NULL,
  `strategy_id` INT UNSIGNED NOT NULL,
  `variant` VARCHAR(100) NOT NULL, -- e.g.: 'Variant A (SMA 50)', 'Variant B (RSI Momentum)'
  `notes` TEXT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_variants_experiment` (`experiment_id`),
  INDEX `idx_variants_strategy` (`strategy_id`),
  CONSTRAINT `fk_variants_experiment`
    FOREIGN KEY (`experiment_id`) REFERENCES `experiments` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_variants_strategy`
    FOREIGN KEY (`strategy_id`) REFERENCES `strategies` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 5. TABLE: IMPLEMENTATIONS (Action Types / Sinks)
-- ============================================================================
CREATE TABLE `implementations` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `class` VARCHAR(100) NOT NULL, -- e.g.: 'SendNotification', 'BrokerOrder'
  `param_list` JSON NOT NULL,
  `type` VARCHAR(50) NOT NULL DEFAULT 'ACTION',
  PRIMARY KEY (`id`),
  INDEX `idx_implementations_class` (`class`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 6. TABLE: ACTIONS (Actions associated with a Variant and Signal)
-- ============================================================================
CREATE TABLE `actions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `variant_id` INT UNSIGNED NOT NULL,
  `implementation_id` INT UNSIGNED NOT NULL,
  `trigger_signal` VARCHAR(10) NOT NULL DEFAULT 'ANY', -- 'BUY', 'SELL', 'HOLD', 'ANY'
  `target_id` BIGINT NOT NULL DEFAULT 0,
  `params` JSON NOT NULL,
  `notes` TEXT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_actions_variant` (`variant_id`),
  INDEX `idx_actions_trigger` (`variant_id`, `trigger_signal`),
  INDEX `idx_actions_implementation` (`implementation_id`),
  CONSTRAINT `fk_actions_variant`
    FOREIGN KEY (`variant_id`) REFERENCES `variants` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_actions_implementation`
    FOREIGN KEY (`implementation_id`) REFERENCES `implementations` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 7. TABLE: ASSET_VARIANTS (Signal Traceability and Historical Log by Asset)
-- ============================================================================
CREATE TABLE `asset_variants` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `variant_id` INT UNSIGNED NOT NULL,
  `experiment_id` INT UNSIGNED NOT NULL,
  `strategy_id` INT UNSIGNED NOT NULL,
  `asset_id` INT UNSIGNED NOT NULL,
  `signal` VARCHAR(10) NOT NULL, -- 'BUY', 'HOLD', 'SELL'
  `price` DECIMAL(18, 4) NOT NULL DEFAULT 0.0000,
  `confidence` DECIMAL(5, 4) NOT NULL DEFAULT 1.0000,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `completed` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_asset_variants_asset_completed` (`asset_id`, `completed`),
  INDEX `idx_asset_variants_exp_asset` (`experiment_id`, `asset_id`, `created_at`),
  INDEX `idx_asset_variants_created_at` (`created_at`),
  CONSTRAINT `fk_av_variant`
    FOREIGN KEY (`variant_id`) REFERENCES `variants` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_av_experiment`
    FOREIGN KEY (`experiment_id`) REFERENCES `experiments` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_av_strategy`
    FOREIGN KEY (`strategy_id`) REFERENCES `strategies` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 8. TABLE: ASSET_ACTIONS (Action/Order Execution and Audit)
-- ============================================================================
CREATE TABLE `asset_actions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `asset_variant_id` INT UNSIGNED NOT NULL,
  `variant_id` INT UNSIGNED NOT NULL,
  `experiment_id` INT UNSIGNED NOT NULL,
  `action_id` INT UNSIGNED NOT NULL,
  `last_attempt_date` DATETIME NULL,
  `attempts` INT UNSIGNED NOT NULL DEFAULT 0,
  `payload` JSON NOT NULL,
  `completed` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_asset_actions_variant_completed` (`asset_variant_id`, `completed`),
  INDEX `idx_asset_actions_exp_action` (`experiment_id`, `action_id`),
  CONSTRAINT `fk_aa_asset_variant`
    FOREIGN KEY (`asset_variant_id`) REFERENCES `asset_variants` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_aa_action`
    FOREIGN KEY (`action_id`) REFERENCES `actions` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- STORED PROCEDURES
-- ============================================================================

DELIMITER $$

-- ----------------------------------------------------------------------------
-- 1. sp_find_experiments
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_experiments`$$
CREATE PROCEDURE `sp_find_experiments`(
  IN p_include_inactive BOOLEAN
)
BEGIN
  SELECT 
    `id`,
    `experiment`,
    `description`,
    `params`,
    `status`,
    `notes`
  FROM `experiments`
  WHERE p_include_inactive = TRUE OR `status` = 'ACTIVE'
  ORDER BY `id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 1.1 sp_get_experiment_by_id
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_get_experiment_by_id`$$
CREATE PROCEDURE `sp_get_experiment_by_id`(
  IN p_experiment_id INT UNSIGNED
)
BEGIN
  SELECT 
    `id`,
    `experiment`,
    `description`,
    `params`,
    `status`,
    `notes`
  FROM `experiments`
  WHERE `id` = p_experiment_id;
END$$

-- ----------------------------------------------------------------------------
-- 1.2 sp_get_experiment_by_name
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_get_experiment_by_name`$$
CREATE PROCEDURE `sp_get_experiment_by_name`(
  IN p_experiment VARCHAR(100)
)
BEGIN
  SELECT 
    `id`,
    `experiment`,
    `description`,
    `params`,
    `status`,
    `notes`
  FROM `experiments`
  WHERE `experiment` = p_experiment;
END$$

-- ----------------------------------------------------------------------------
-- 2. sp_find_experiment_assets
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_experiment_assets`$$
CREATE PROCEDURE `sp_find_experiment_assets`(
  IN p_experiment_id INT UNSIGNED
)
BEGIN
  SELECT 
    `id`,
    `experiment_id`,
    `asset_id`,
    `status`,
    `created_at`
  FROM `experiment_assets`
  WHERE `experiment_id` = p_experiment_id AND `status` = 'ACTIVE'
  ORDER BY `asset_id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 3. sp_find_experiments_by_asset
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_experiments_by_asset`$$
CREATE PROCEDURE `sp_find_experiments_by_asset`(
  IN p_asset_id INT UNSIGNED
)
BEGIN
  SELECT 
    e.`id`,
    e.`experiment`,
    e.`description`,
    e.`params`,
    e.`status`,
    e.`notes`
  FROM `experiments` e
  INNER JOIN `experiment_assets` ea ON ea.`experiment_id` = e.`id`
  WHERE ea.`asset_id` = p_asset_id
    AND ea.`status` = 'ACTIVE'
    AND e.`status` = 'ACTIVE'
  ORDER BY e.`id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 4. sp_find_variants_by_experiment
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_variants_by_experiment`$$
CREATE PROCEDURE `sp_find_variants_by_experiment`(
  IN p_experiment_id INT UNSIGNED
)
BEGIN
  SELECT 
    v.`id`,
    v.`experiment_id`,
    v.`strategy_id`,
    v.`variant`,
    v.`notes`,
    st.`name` AS `strategy_name`,
    st.`class`,
    st.`params`,
    st.`status` AS `strategy_status`
  FROM `variants` v
  INNER JOIN `strategies` st ON st.`id` = v.`strategy_id`
  WHERE (p_experiment_id = 0 OR v.`experiment_id` = p_experiment_id)
    AND st.`status` = 'ACTIVE'
  ORDER BY v.`id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 5. sp_start_asset_variant
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_start_asset_variant`$$
CREATE PROCEDURE `sp_start_asset_variant`(
  IN p_asset_id INT UNSIGNED,
  IN p_variant_id INT UNSIGNED,
  IN p_signal VARCHAR(10),
  IN p_price DECIMAL(18, 4)
)
BEGIN
  DECLARE v_asset_variant_id INT UNSIGNED;
  DECLARE v_experiment_id INT UNSIGNED;
  DECLARE v_strategy_id INT UNSIGNED;

  SELECT `experiment_id`, `strategy_id`
  INTO v_experiment_id, v_strategy_id
  FROM `variants`
  WHERE `id` = p_variant_id;

  INSERT INTO `asset_variants` (
    `variant_id`,
    `experiment_id`,
    `strategy_id`,
    `asset_id`,
    `signal`,
    `price`,
    `created_at`
  ) VALUES (
    p_variant_id,
    v_experiment_id,
    v_strategy_id,
    p_asset_id,
    p_signal,
    p_price,
    NOW()
  );

  SET v_asset_variant_id = LAST_INSERT_ID();

  INSERT INTO `asset_actions` (
    `asset_variant_id`,
    `variant_id`,
    `experiment_id`,
    `action_id`,
    `payload`,
    `created_at`
  )
  SELECT 
    v_asset_variant_id,
    p_variant_id,
    v_experiment_id,
    a.`id`,
    a.`params`,
    NOW()
  FROM `actions` a
  WHERE a.`variant_id` = p_variant_id
    AND (a.`trigger_signal` = 'ANY' OR a.`trigger_signal` = p_signal);

  SELECT 
    `id`,
    `variant_id`,
    `experiment_id`,
    `strategy_id`,
    `asset_id`,
    `signal`,
    `price`,
    `created_at`,
    `completed`
  FROM `asset_variants`
  WHERE `id` = v_asset_variant_id;
END$$

-- ----------------------------------------------------------------------------
-- 6. sp_complete_asset_variant
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_complete_asset_variant`$$
CREATE PROCEDURE `sp_complete_asset_variant`(
  IN p_asset_variant_id INT UNSIGNED
)
BEGIN
  UPDATE `asset_variants`
  SET `completed` = NOW()
  WHERE `id` = p_asset_variant_id;

  SELECT 
    `id`,
    `variant_id`,
    `experiment_id`,
    `strategy_id`,
    `asset_id`,
    `signal`,
    `price`,
    `created_at`,
    `completed`
  FROM `asset_variants`
  WHERE `id` = p_asset_variant_id;
END$$

-- ----------------------------------------------------------------------------
-- 7. sp_find_asset_variants
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_asset_variants`$$
CREATE PROCEDURE `sp_find_asset_variants`(
  IN p_asset_id INT UNSIGNED,
  IN p_from_date DATETIME,
  IN p_to_date DATETIME,
  IN p_include_completed BOOLEAN
)
BEGIN
  SELECT 
    `id`,
    `variant_id`,
    `experiment_id`,
    `strategy_id`,
    `asset_id`,
    `signal`,
    `price`,
    `created_at`,
    `completed`
  FROM `asset_variants`
  WHERE `asset_id` = p_asset_id
    AND (p_from_date IS NULL OR `created_at` >= p_from_date)
    AND (p_to_date IS NULL OR `created_at` <= p_to_date)
    AND (p_include_completed = TRUE OR `completed` IS NULL)
  ORDER BY `created_at` DESC;
END$$

-- ----------------------------------------------------------------------------
-- 8. sp_find_asset_actions
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_asset_actions`$$
CREATE PROCEDURE `sp_find_asset_actions`(
  IN p_asset_variant_id INT UNSIGNED,
  IN p_include_completed BOOLEAN
)
BEGIN
  SELECT 
    aa.`id` AS `assetaction.id`,
    aa.`asset_variant_id` AS `assetaction.asset_variant_id`,
    aa.`variant_id` AS `assetaction.variant_id`,
    aa.`experiment_id` AS `assetaction.experiment_id`,
    aa.`action_id` AS `assetaction.action_id`,
    aa.`last_attempt_date` AS `assetaction.last_attempt_date`,
    aa.`attempts` AS `assetaction.attempts`,
    aa.`payload` AS `assetaction.payload`,
    aa.`completed` AS `assetaction.completed`,
    i.`id` AS `implementation.id`,
    i.`class` AS `implementation.class`,
    i.`param_list` AS `implementation.param_list`,
    i.`type` AS `implementation.type`
  FROM `asset_actions` aa
  INNER JOIN `actions` a ON a.`id` = aa.`action_id`
  INNER JOIN `implementations` i ON i.`id` = a.`implementation_id`
  WHERE aa.`asset_variant_id` = p_asset_variant_id
    AND (p_include_completed = TRUE OR aa.`completed` IS NULL)
  ORDER BY aa.`id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 9. sp_start_asset_action_attempt
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_start_asset_action_attempt`$$
CREATE PROCEDURE `sp_start_asset_action_attempt`(
  IN p_asset_action_id INT UNSIGNED
)
BEGIN
  UPDATE `asset_actions`
  SET 
    `attempts` = `attempts` + 1,
    `last_attempt_date` = NOW()
  WHERE `id` = p_asset_action_id;

  SELECT 
    `id`,
    `asset_variant_id`,
    `variant_id`,
    `experiment_id`,
    `action_id`,
    `last_attempt_date`,
    `attempts`,
    `payload`,
    `completed`
  FROM `asset_actions`
  WHERE `id` = p_asset_action_id;
END$$

-- ----------------------------------------------------------------------------
-- 10. sp_complete_asset_action
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_complete_asset_action`$$
CREATE PROCEDURE `sp_complete_asset_action`(
  IN p_asset_action_id INT UNSIGNED
)
BEGIN
  UPDATE `asset_actions`
  SET `completed` = NOW()
  WHERE `id` = p_asset_action_id;

  SELECT 
    `id`,
    `asset_variant_id`,
    `variant_id`,
    `experiment_id`,
    `action_id`,
    `last_attempt_date`,
    `attempts`,
    `payload`,
    `completed`
  FROM `asset_actions`
  WHERE `id` = p_asset_action_id;
END$$

-- ----------------------------------------------------------------------------
-- 11. sp_find_actions_by_variant
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_actions_by_variant`$$
CREATE PROCEDURE `sp_find_actions_by_variant`(
  IN p_variant_id INT UNSIGNED
)
BEGIN
  SELECT 
    a.`id` AS `action.id`,
    a.`variant_id` AS `action.variant_id`,
    a.`implementation_id` AS `action.implementation_id`,
    a.`trigger_signal` AS `action.trigger_signal`,
    a.`target_id` AS `action.target_id`,
    a.`params` AS `action.params`,
    a.`notes` AS `action.notes`,
    i.`id` AS `implementation.id`,
    i.`class` AS `implementation.class`,
    i.`param_list` AS `implementation.param_list`,
    i.`type` AS `implementation.type`
  FROM `actions` a
  INNER JOIN `implementations` i ON i.`id` = a.`implementation_id`
  WHERE (p_variant_id = 0 OR a.`variant_id` = p_variant_id)
  ORDER BY a.`id` ASC;
END$$

DELIMITER ;

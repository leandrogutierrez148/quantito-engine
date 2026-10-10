-- ============================================================================
-- QUANTITO ENGINE - MySQL Database Schema & Stored Procedures
-- Engine: InnoDB | Character Set: utf8mb4 | Collation: utf8mb4_unicode_ci
-- ============================================================================

CREATE DATABASE IF NOT EXISTS `quantito-engine`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `quantito-engine`;

SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `signals`;
DROP TABLE IF EXISTS `execution_variants`;
DROP TABLE IF EXISTS `executions`;
DROP TABLE IF EXISTS `variants`;
DROP TABLE IF EXISTS `strategies`;
DROP TABLE IF EXISTS `experiment_assets`;
DROP TABLE IF EXISTS `experiments`;
SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 1. TABLE: EXPERIMENTS (Quantitative Experiments Specification)
-- ============================================================================
CREATE TABLE `experiments` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `description` VARCHAR(255) NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- 'DRAFT', 'ACTIVE', 'PAUSED', 'COMPLETED'
  `execution_mode` VARCHAR(20) NOT NULL DEFAULT 'PAPER', -- 'BACKTEST', 'PAPER', 'LIVE'
  `start_date` DATETIME NULL,
  `end_date` DATETIME NULL,
  `frequency` VARCHAR(50) NOT NULL DEFAULT '1d', -- e.g. '15m', '1h', '1d'
  `max_capital` DECIMAL(18, 4) NOT NULL DEFAULT 10000.0000,
  `currency` VARCHAR(10) NOT NULL DEFAULT 'USD',
  `assets` JSON NOT NULL, -- JSON Array of assets, e.g. ["AAPL", "MSFT", "NVDA"]
  `params` JSON NULL, -- Optional execution settings (e.g. benchmark, slippage, commission)
  `notes` TEXT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_experiments_status` (`status`),
  INDEX `idx_experiments_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 2. TABLE: STRATEGIES (Quantitative Strategies Catalog & Param Specs)
-- ============================================================================
CREATE TABLE `strategies` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `class` VARCHAR(100) NOT NULL, -- e.g.: 'fifty_day_avg', 'price_momentum'
  `param_definitions` JSON NOT NULL, -- Parameter specifications (types, defaults, required)
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_strategies_class` (`class`),
  INDEX `idx_strategies_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 3. TABLE: VARIANTS (Experiment Variants / Parametrized Strategy Setups)
-- ============================================================================
CREATE TABLE `variants` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `experiment_id` INT UNSIGNED NOT NULL,
  `strategy_id` INT UNSIGNED NOT NULL,
  `name` VARCHAR(100) NOT NULL, -- e.g.: 'Variant A (SMA 50)', 'Variant B (Price Momentum)'
  `params` JSON NOT NULL, -- Concrete runtime parameters provided to the strategy
  `capital_allocation_pct` DECIMAL(5, 4) NOT NULL DEFAULT 1.0000,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
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
-- 4. TABLE: EXECUTIONS (Discrete Experiment Run Lifecycle)
-- ============================================================================
CREATE TABLE `executions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `experiment_id` INT UNSIGNED NOT NULL,
  `trigger_type` VARCHAR(20) NOT NULL DEFAULT 'MANUAL', -- 'MANUAL', 'SCHEDULED', 'BACKTEST', 'WEBHOOK'
  `triggered_by` VARCHAR(100) NOT NULL DEFAULT 'system',
  `as_of_time` DATETIME NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- 'PENDING', 'RUNNING', 'COMPLETED', 'PARTIAL_SUCCESS', 'FAILED'
  `started_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `completed_at` DATETIME NULL,
  `error_message` TEXT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_executions_experiment` (`experiment_id`, `as_of_time`),
  INDEX `idx_executions_status` (`status`),
  CONSTRAINT `fk_executions_experiment`
    FOREIGN KEY (`experiment_id`) REFERENCES `experiments` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 5. TABLE: EXECUTION_VARIANTS (Variant Execution Outcome per Run)
-- ============================================================================
CREATE TABLE `execution_variants` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `execution_id` INT UNSIGNED NOT NULL,
  `variant_id` INT UNSIGNED NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- 'PENDING', 'RUNNING', 'COMPLETED', 'FAILED'
  `signals_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `duration_ms` INT UNSIGNED NOT NULL DEFAULT 0,
  `error_message` TEXT NULL,
  `started_at` DATETIME NULL,
  `completed_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_execution_variant` (`execution_id`, `variant_id`),
  INDEX `idx_exec_variants_variant` (`variant_id`),
  INDEX `idx_exec_variants_status` (`status`),
  CONSTRAINT `fk_exec_variants_execution`
    FOREIGN KEY (`execution_id`) REFERENCES `executions` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_exec_variants_variant`
    FOREIGN KEY (`variant_id`) REFERENCES `variants` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- 6. TABLE: SIGNALS (Signal Traceability Log linked to Execution Variant)
-- ============================================================================
CREATE TABLE `signals` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `execution_variant_id` INT UNSIGNED NOT NULL,
  `asset` VARCHAR(20) NOT NULL,
  `signal` VARCHAR(10) NOT NULL, -- 'BUY', 'HOLD', 'SELL'
  `price` DECIMAL(18, 4) NOT NULL DEFAULT 0.0000,
  `confidence` DECIMAL(5, 4) NOT NULL DEFAULT 1.0000,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_signals_exec_variant` (`execution_variant_id`),
  INDEX `idx_signals_asset` (`asset`, `created_at`),
  INDEX `idx_signals_signal` (`signal`, `created_at`),
  CONSTRAINT `fk_signals_execution_variant`
    FOREIGN KEY (`execution_variant_id`) REFERENCES `execution_variants` (`id`)
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
    `name`,
    `description`,
    `status`,
    `execution_mode`,
    `start_date`,
    `end_date`,
    `frequency`,
    `max_capital`,
    `currency`,
    `assets`,
    `params`,
    `notes`,
    `created_at`,
    `updated_at`
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
    `name`,
    `description`,
    `status`,
    `execution_mode`,
    `start_date`,
    `end_date`,
    `frequency`,
    `max_capital`,
    `currency`,
    `assets`,
    `params`,
    `notes`,
    `created_at`,
    `updated_at`
  FROM `experiments`
  WHERE `id` = p_experiment_id;
END$$

-- ----------------------------------------------------------------------------
-- 1.2 sp_get_experiment_by_name
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_get_experiment_by_name`$$
CREATE PROCEDURE `sp_get_experiment_by_name`(
  IN p_name VARCHAR(100)
)
BEGIN
  SELECT 
    `id`,
    `name`,
    `description`,
    `status`,
    `execution_mode`,
    `start_date`,
    `end_date`,
    `frequency`,
    `max_capital`,
    `currency`,
    `assets`,
    `params`,
    `notes`,
    `created_at`,
    `updated_at`
  FROM `experiments`
  WHERE `name` = p_name;
END$$

-- ----------------------------------------------------------------------------
-- 2. sp_find_variants_by_experiment
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
    v.`name`,
    v.`params`,
    v.`capital_allocation_pct`,
    v.`status`,
    v.`notes`,
    st.`name` AS `strategy_name`,
    st.`class`,
    st.`param_definitions`,
    st.`status` AS `strategy_status`
  FROM `variants` v
  INNER JOIN `strategies` st ON st.`id` = v.`strategy_id`
  WHERE (p_experiment_id = 0 OR v.`experiment_id` = p_experiment_id)
    AND st.`status` = 'ACTIVE'
    AND v.`status` = 'ACTIVE'
  ORDER BY v.`id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 3. sp_create_execution
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_create_execution`$$
CREATE PROCEDURE `sp_create_execution`(
  IN p_experiment_id INT UNSIGNED,
  IN p_trigger_type VARCHAR(20),
  IN p_triggered_by VARCHAR(100),
  IN p_as_of_time DATETIME
)
BEGIN
  DECLARE v_execution_id INT UNSIGNED;

  INSERT INTO `executions` (
    `experiment_id`,
    `trigger_type`,
    `triggered_by`,
    `as_of_time`,
    `status`,
    `started_at`
  ) VALUES (
    p_experiment_id,
    p_trigger_type,
    p_triggered_by,
    p_as_of_time,
    'RUNNING',
    NOW()
  );

  SET v_execution_id = LAST_INSERT_ID();

  SELECT 
    `id`,
    `experiment_id`,
    `trigger_type`,
    `triggered_by`,
    `as_of_time`,
    `status`,
    `started_at`,
    `completed_at`,
    `error_message`,
    `created_at`,
    `updated_at`
  FROM `executions`
  WHERE `id` = v_execution_id;
END$$

-- ----------------------------------------------------------------------------
-- 4. sp_update_execution
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_update_execution`$$
CREATE PROCEDURE `sp_update_execution`(
  IN p_execution_id INT UNSIGNED,
  IN p_status VARCHAR(20),
  IN p_error_message TEXT
)
BEGIN
  UPDATE `executions`
  SET 
    `status` = p_status,
    `completed_at` = NOW(),
    `error_message` = p_error_message
  WHERE `id` = p_execution_id;

  SELECT 
    `id`,
    `experiment_id`,
    `trigger_type`,
    `triggered_by`,
    `as_of_time`,
    `status`,
    `started_at`,
    `completed_at`,
    `error_message`,
    `created_at`,
    `updated_at`
  FROM `executions`
  WHERE `id` = p_execution_id;
END$$

-- ----------------------------------------------------------------------------
-- 5. sp_create_execution_variant
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_create_execution_variant`$$
CREATE PROCEDURE `sp_create_execution_variant`(
  IN p_execution_id INT UNSIGNED,
  IN p_variant_id INT UNSIGNED
)
BEGIN
  DECLARE v_id INT UNSIGNED;

  INSERT INTO `execution_variants` (
    `execution_id`,
    `variant_id`,
    `status`,
    `started_at`
  ) VALUES (
    p_execution_id,
    p_variant_id,
    'RUNNING',
    NOW()
  );

  SET v_id = LAST_INSERT_ID();

  SELECT 
    `id`,
    `execution_id`,
    `variant_id`,
    `status`,
    `signals_count`,
    `duration_ms`,
    `error_message`,
    `started_at`,
    `completed_at`,
    `created_at`,
    `updated_at`
  FROM `execution_variants`
  WHERE `id` = v_id;
END$$

-- ----------------------------------------------------------------------------
-- 6. sp_update_execution_variant
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_update_execution_variant`$$
CREATE PROCEDURE `sp_update_execution_variant`(
  IN p_id INT UNSIGNED,
  IN p_status VARCHAR(20),
  IN p_signals_count INT UNSIGNED,
  IN p_duration_ms INT UNSIGNED,
  IN p_error_message TEXT
)
BEGIN
  UPDATE `execution_variants`
  SET 
    `status` = p_status,
    `signals_count` = p_signals_count,
    `duration_ms` = p_duration_ms,
    `error_message` = p_error_message,
    `completed_at` = NOW()
  WHERE `id` = p_id;

  SELECT 
    `id`,
    `execution_id`,
    `variant_id`,
    `status`,
    `signals_count`,
    `duration_ms`,
    `error_message`,
    `started_at`,
    `completed_at`,
    `created_at`,
    `updated_at`
  FROM `execution_variants`
  WHERE `id` = p_id;
END$$

-- ----------------------------------------------------------------------------
-- 7. sp_save_signal
-- Saves signal outcome linked to the execution variant and asset.
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_save_signal`$$
CREATE PROCEDURE `sp_save_signal`(
  IN p_execution_variant_id INT UNSIGNED,
  IN p_asset VARCHAR(20),
  IN p_signal VARCHAR(10),
  IN p_price DECIMAL(18, 4),
  IN p_confidence DECIMAL(5, 4)
)
BEGIN
  DECLARE v_signal_id INT UNSIGNED;

  INSERT INTO `signals` (
    `execution_variant_id`,
    `asset`,
    `signal`,
    `price`,
    `confidence`,
    `created_at`
  ) VALUES (
    p_execution_variant_id,
    p_asset,
    p_signal,
    p_price,
    p_confidence,
    NOW()
  );

  SET v_signal_id = LAST_INSERT_ID();

  SELECT 
    `id`,
    `execution_variant_id`,
    `asset`,
    `signal`,
    `price`,
    `confidence`,
    `created_at`
  FROM `signals`
  WHERE `id` = v_signal_id;
END$$

-- ----------------------------------------------------------------------------
-- 8. sp_find_signals_by_execution
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_signals_by_execution`$$
CREATE PROCEDURE `sp_find_signals_by_execution`(
  IN p_execution_id INT UNSIGNED
)
BEGIN
  SELECT 
    s.`id`,
    s.`execution_variant_id`,
    ev.`execution_id`,
    ev.`variant_id`,
    v.`name` AS `variant_name`,
    s.`asset`,
    s.`signal`,
    s.`price`,
    s.`confidence`,
    s.`created_at`
  FROM `signals` s
  INNER JOIN `execution_variants` ev ON ev.`id` = s.`execution_variant_id`
  INNER JOIN `variants` v ON v.`id` = ev.`variant_id`
  WHERE ev.`execution_id` = p_execution_id
  ORDER BY s.`id` ASC;
END$$

-- ----------------------------------------------------------------------------
-- 9. sp_find_signals_by_asset
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_find_signals_by_asset`$$
CREATE PROCEDURE `sp_find_signals_by_asset`(
  IN p_asset VARCHAR(20),
  IN p_from_date DATETIME,
  IN p_to_date DATETIME
)
BEGIN
  SELECT 
    s.`id`,
    s.`execution_variant_id`,
    ev.`execution_id`,
    ev.`variant_id`,
    v.`name` AS `variant_name`,
    s.`asset`,
    s.`signal`,
    s.`price`,
    s.`confidence`,
    s.`created_at`
  FROM `signals` s
  INNER JOIN `execution_variants` ev ON ev.`id` = s.`execution_variant_id`
  INNER JOIN `variants` v ON v.`id` = ev.`variant_id`
  WHERE s.`asset` = p_asset
    AND (p_from_date IS NULL OR s.`created_at` >= p_from_date)
    AND (p_to_date IS NULL OR s.`created_at` <= p_to_date)
  ORDER BY s.`created_at` DESC;
END$$

DELIMITER ;
